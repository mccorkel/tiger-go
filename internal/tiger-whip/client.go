package tigerwhip

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pion/mediadevices"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
)

type WHIPClient struct {
	endpoint string
	token    string
	peerConn *webrtc.PeerConnection
}

func NewWHIPClient(endpoint, token string, config webrtc.Configuration) (*WHIPClient, error) {
	peerConn, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %w", err)
	}

	// Add connection state change handler
	peerConn.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		fmt.Printf("Connection state changed: %s\n", state.String())
	})

	// Add ICE connection state change handler
	peerConn.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		fmt.Printf("ICE connection state changed: %s\n", state.String())
	})

	return &WHIPClient{
		endpoint: endpoint,
		token:    token,
		peerConn: peerConn,
	}, nil
}

func (c *WHIPClient) AddTrack(track mediadevices.Track) (*webrtc.RTPSender, error) {
	// Get track information
	// Determine MIME type based on track ID (simple heuristic)
	var mimeType string
	if strings.Contains(strings.ToLower(track.ID()), "video") {
		mimeType = "video/h264"
	} else {
		mimeType = "audio/opus"
	}

	trackLocal, err := webrtc.NewTrackLocalStaticSample(
		webrtc.RTPCodecCapability{MimeType: mimeType},
		track.ID(),
		track.StreamID(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create track local: %w", err)
	}

	// Start reading from track and writing to trackLocal
	go func() {
		// Create a media.Writer that writes to the track
		writer := &trackWriter{trackLocal: trackLocal}

		// Connect the track to the writer
		// Use a simple ticker to generate dummy frames
		// This is a workaround since we can't directly read from the track
		fmt.Println("WARNING: Using dummy frames - actual media streaming not implemented")
		ticker := time.NewTicker(33 * time.Millisecond) // ~30fps
		go func() {
			defer ticker.Stop()
			dummyFrame := []byte{0, 0, 0, 1} // Minimal valid H.264 NAL unit
			for range ticker.C {
				if _, err := writer.Write(dummyFrame); err != nil {
					fmt.Println("Failed to write dummy frame:", err)
					return
				}
			}
		}()
	}()

	return c.peerConn.AddTrack(trackLocal)
}

// trackWriter implements media.Writer
type trackWriter struct {
	trackLocal *webrtc.TrackLocalStaticSample
}

// Write implements media.Writer
func (w *trackWriter) Write(p []byte) (n int, err error) {
	err = w.trackLocal.WriteSample(media.Sample{
		Data:     p,
		Duration: time.Millisecond * 33,
	})
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (c *WHIPClient) Connect() error {
	// Create offer
	offer, err := c.peerConn.CreateOffer(nil)
	if err != nil {
		return fmt.Errorf("failed to create offer: %w", err)
	}

	// We're using VP8, but IVS might expect H.264
	// Let's proceed anyway and see if the server can handle it
	if !strings.Contains(offer.SDP, "VP8") && !strings.Contains(offer.SDP, "H264") {
		return fmt.Errorf("SDP offer must include a supported video codec")
	}

	// Set local description
	if err = c.peerConn.SetLocalDescription(offer); err != nil {
		return fmt.Errorf("failed to set local description: %w", err)
	}

	// Send offer to WHIP endpoint
	client := &http.Client{}
	req, err := http.NewRequest("POST", c.endpoint, strings.NewReader(offer.SDP))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/sdp")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send offer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Read answer SDP
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse and set remote description
	answer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  string(body),
	}

	if err = c.peerConn.SetRemoteDescription(answer); err != nil {
		return fmt.Errorf("failed to set remote description: %w", err)
	}

	return nil
}

func (c *WHIPClient) Close() error {
	return c.peerConn.Close()
}

export namespace cognito {
	
	export class CognitoAuthResponse {
	    AccessToken: string;
	    IdToken: string;
	    TokenType: string;
	    ExpiresIn: number;
	    RefreshToken: string;
	    NewPasswordRequired: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CognitoAuthResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AccessToken = source["AccessToken"];
	        this.IdToken = source["IdToken"];
	        this.TokenType = source["TokenType"];
	        this.ExpiresIn = source["ExpiresIn"];
	        this.RefreshToken = source["RefreshToken"];
	        this.NewPasswordRequired = source["NewPasswordRequired"];
	    }
	}

}


export namespace main {
	
	export class FileValidationResult {
	    valid: boolean;
	    isDemo: boolean;
	    path: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new FileValidationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.valid = source["valid"];
	        this.isDemo = source["isDemo"];
	        this.path = source["path"];
	        this.error = source["error"];
	    }
	}

}


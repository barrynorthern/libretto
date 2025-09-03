export namespace main {
	
	export class SceneDTO {
	    id: string;
	    title: string;
	    summary: string;
	    content: string;
	    created: string;
	
	    static createFrom(source: any = {}) {
	        return new SceneDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.summary = source["summary"];
	        this.content = source["content"];
	        this.created = source["created"];
	    }
	}

}


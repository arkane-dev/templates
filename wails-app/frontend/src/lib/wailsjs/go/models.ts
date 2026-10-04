export namespace main {
	
	export class AppInfo {
	    name: string;
	    version: string;
	    goVersion: string;
	    os: string;
	    arch: string;
	    hostname: string;
	    cpus: number;
	
	    static createFrom(source: any = {}) {
	        return new AppInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	        this.goVersion = source["goVersion"];
	        this.os = source["os"];
	        this.arch = source["arch"];
	        this.hostname = source["hostname"];
	        this.cpus = source["cpus"];
	    }
	}
	export class Settings {
	    accent: string;
	    glowSize: number;
	    scanlines: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.accent = source["accent"];
	        this.glowSize = source["glowSize"];
	        this.scanlines = source["scanlines"];
	    }
	}

}


export namespace liveroom {
	
	export class ArtQuality {
	    default: boolean;
	    html: string;
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new ArtQuality(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.default = source["default"];
	        this.html = source["html"];
	        this.url = source["url"];
	    }
	}
	export class Quality {
	    name: string;
	    url: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new Quality(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.url = source["url"];
	        this.type = source["type"];
	    }
	}
	export class LiveRoom {
	    platform: string;
	    platformName: string;
	    roomId: string;
	    roomName: string;
	    anchor: string;
	    avatar: string;
	    onLineCount: number;
	    liveUrl: string;
	    quality: Quality[];
	    isLive: boolean;
	    screenshot: string;
	    gameFullName: string;
	
	    static createFrom(source: any = {}) {
	        return new LiveRoom(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.platform = source["platform"];
	        this.platformName = source["platformName"];
	        this.roomId = source["roomId"];
	        this.roomName = source["roomName"];
	        this.anchor = source["anchor"];
	        this.avatar = source["avatar"];
	        this.onLineCount = source["onLineCount"];
	        this.liveUrl = source["liveUrl"];
	        this.quality = this.convertValues(source["quality"], Quality);
	        this.isLive = source["isLive"];
	        this.screenshot = source["screenshot"];
	        this.gameFullName = source["gameFullName"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LiveRoomInfo {
	    platform: string;
	    platformName: string;
	    roomId: string;
	    roomName: string;
	    anchor: string;
	    avatar: string;
	    onLineCount: number;
	    screenshot: string;
	    gameFullName: string;
	    liveStatus: number;
	
	    static createFrom(source: any = {}) {
	        return new LiveRoomInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.platform = source["platform"];
	        this.platformName = source["platformName"];
	        this.roomId = source["roomId"];
	        this.roomName = source["roomName"];
	        this.anchor = source["anchor"];
	        this.avatar = source["avatar"];
	        this.onLineCount = source["onLineCount"];
	        this.screenshot = source["screenshot"];
	        this.gameFullName = source["gameFullName"];
	        this.liveStatus = source["liveStatus"];
	    }
	}

}


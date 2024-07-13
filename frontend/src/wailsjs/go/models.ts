export namespace main {
	
	export class Config {
	    theme?: string;
	    useSystemTitleBar?: boolean;
	    enableLogging?: boolean;
	    enableTrace?: boolean;
	    enableDebug?: boolean;
	    enableInfo?: boolean;
	    enableWarn?: boolean;
	    enableError?: boolean;
	    enableFatal?: boolean;
	    maxLogFiles?: number;
	    language?: string;
	    saveWindowStatus?: boolean;
	    windowStartState?: number;
	    windowStartPositionX?: number;
	    windowStartPositionY?: number;
	    windowStartSizeX?: number;
	    windowStartSizeY?: number;
	    windowScale?: number;
	    opacity?: number;
	    windowEffect?: number;
	    checkForUpdates?: boolean;
	    lastUpdateCheck?: number;
	    folderNamePattern?: string;
	    wordFileNamePattern?: string;
	    fileNamePattern?: string;
	    ilCellName?: string;
	    ilceCellName?: string;
	    mahalleCellName?: string;
	    adaCellName?: string;
	    parselCellName?: string;
	    alanCellName?: string;
	    paftaCellName?: string;
	    parselSorguHeadless?: boolean;
	    ciltCellName?: string;
	    sayfaCellName?: string;
	    tapuNamePattern?: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.theme = source["theme"];
	        this.useSystemTitleBar = source["useSystemTitleBar"];
	        this.enableLogging = source["enableLogging"];
	        this.enableTrace = source["enableTrace"];
	        this.enableDebug = source["enableDebug"];
	        this.enableInfo = source["enableInfo"];
	        this.enableWarn = source["enableWarn"];
	        this.enableError = source["enableError"];
	        this.enableFatal = source["enableFatal"];
	        this.maxLogFiles = source["maxLogFiles"];
	        this.language = source["language"];
	        this.saveWindowStatus = source["saveWindowStatus"];
	        this.windowStartState = source["windowStartState"];
	        this.windowStartPositionX = source["windowStartPositionX"];
	        this.windowStartPositionY = source["windowStartPositionY"];
	        this.windowStartSizeX = source["windowStartSizeX"];
	        this.windowStartSizeY = source["windowStartSizeY"];
	        this.windowScale = source["windowScale"];
	        this.opacity = source["opacity"];
	        this.windowEffect = source["windowEffect"];
	        this.checkForUpdates = source["checkForUpdates"];
	        this.lastUpdateCheck = source["lastUpdateCheck"];
	        this.folderNamePattern = source["folderNamePattern"];
	        this.wordFileNamePattern = source["wordFileNamePattern"];
	        this.fileNamePattern = source["fileNamePattern"];
	        this.ilCellName = source["ilCellName"];
	        this.ilceCellName = source["ilceCellName"];
	        this.mahalleCellName = source["mahalleCellName"];
	        this.adaCellName = source["adaCellName"];
	        this.parselCellName = source["parselCellName"];
	        this.alanCellName = source["alanCellName"];
	        this.paftaCellName = source["paftaCellName"];
	        this.parselSorguHeadless = source["parselSorguHeadless"];
	        this.ciltCellName = source["ciltCellName"];
	        this.sayfaCellName = source["sayfaCellName"];
	        this.tapuNamePattern = source["tapuNamePattern"];
	    }
	}
	export class Properties {
	    ParselNo: string;
	    Alan: string;
	    Mevkii: string;
	    Nitelik: string;
	    Ada: string;
	    Il: string;
	    Ilce: string;
	    Pafta: string;
	    Mahalle: string;
	
	    static createFrom(source: any = {}) {
	        return new Properties(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ParselNo = source["ParselNo"];
	        this.Alan = source["Alan"];
	        this.Mevkii = source["Mevkii"];
	        this.Nitelik = source["Nitelik"];
	        this.Ada = source["Ada"];
	        this.Il = source["Il"];
	        this.Ilce = source["Ilce"];
	        this.Pafta = source["Pafta"];
	        this.Mahalle = source["Mahalle"];
	    }
	}
	export class QueryParams {
	    province: string;
	    district: string;
	    neighborhood: string;
	    block: string;
	    parcel: string;
	
	    static createFrom(source: any = {}) {
	        return new QueryParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.province = source["province"];
	        this.district = source["district"];
	        this.neighborhood = source["neighborhood"];
	        this.block = source["block"];
	        this.parcel = source["parcel"];
	    }
	}
	export class Tapu {
	
	
	    static createFrom(source: any = {}) {
	        return new Tapu(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class UpdateInfo {
	    updateAvailable: boolean;
	    currentVersion: string;
	    latestVersion: string;
	    name: string;
	    releaseNotes: string;
	    downloadUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.updateAvailable = source["updateAvailable"];
	        this.currentVersion = source["currentVersion"];
	        this.latestVersion = source["latestVersion"];
	        this.name = source["name"];
	        this.releaseNotes = source["releaseNotes"];
	        this.downloadUrl = source["downloadUrl"];
	    }
	}

}


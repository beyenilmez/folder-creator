// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';

export function AddParselSorguFields(arg1:string,arg2:string,arg3:string,arg4:string,arg5:string,arg6:string,arg7:string,arg8:string,arg9:boolean):Promise<void>;

export function AddTapuToExcel(arg1:string,arg2:string,arg3:string,arg4:string,arg5:string):Promise<string>;

export function CheckForUpdate():Promise<main.UpdateInfo>;

export function CreateFolders(arg1:string,arg2:string,arg3:string,arg4:string,arg5:string,arg6:string,arg7:string,arg8:string):Promise<string>;

export function CreateFoldersV2(arg1:string,arg2:string,arg3:string):Promise<string>;

export function GetConfig():Promise<main.Config>;

export function GetConfigField(arg1:string):Promise<any>;

export function GetCopyFolderDialog():Promise<string>;

export function GetExcelFileDialog():Promise<string>;

export function GetFileDialog():Promise<string>;

export function GetLoadConfigPath():Promise<string>;

export function GetTargetFolderDialog():Promise<string>;

export function GetVersion():Promise<string>;

export function GetWordFileDialog():Promise<string>;

export function InitParselSorgu(arg1:boolean):Promise<void>;

export function NeedsAdminPrivileges():Promise<boolean>;

export function OpenFile(arg1:string):Promise<void>;

export function OpenFileInExplorer(arg1:string):Promise<void>;

export function ParseTapu(arg1:string):Promise<main.Tapu>;

export function ParselSorgu(arg1:main.QueryParams):Promise<main.Properties>;

export function ReadConfig(arg1:string):Promise<void>;

export function RestartApplication(arg1:boolean,arg2:Array<string>):Promise<void>;

export function SaveConfigDialog():Promise<void>;

export function SendNotification(arg1:string,arg2:string,arg3:string,arg4:string):Promise<void>;

export function SetConfigField(arg1:string,arg2:any):Promise<void>;

export function Update(arg1:string):Promise<void>;

export function UpdateAsAdmin(arg1:string):Promise<void>;

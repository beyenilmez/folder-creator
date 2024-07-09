import { useState } from "react";
import { Button } from "./ui/button";
import {
  CreateFolders,
  GetCopyFolderDialog,
  GetExcelFileDialog,
  GetTargetFolderDialog,
  SendNotification,
} from "@/wailsjs/go/main/App";

export function Home() {
  const [excelPath, setExcelPath] = useState<string>("");
  const [copyFolder, setCopyFolder] = useState<string>("");
  const [targetFolder, setTargetFolder] = useState<string>("");

  const handleExcelFileDialog = () => {
    GetExcelFileDialog().then((path) => {
      setExcelPath(path);
    });
  };

  const handleCopyFolder = () => {
    GetCopyFolderDialog().then((path) => {
      setCopyFolder(path);
    });
  };

  const handleTargetFolder = () => {
    GetTargetFolderDialog().then((path) => {
      setTargetFolder(path);
    });
  };

  return (
    <div className="flex flex-col justify-center items-center gap-12 w-full h-full">
      <div className="flex flex-col gap-2 text-center">
        <Button variant={"outline"} onClick={handleExcelFileDialog}>
          Excel Dosyası Seçin
        </Button>
        <label>{excelPath ? excelPath : "Dosya seçilmedi..."}</label>
      </div>
      <div className="flex flex-col gap-2 text-center">
        <Button variant={"outline"} onClick={handleCopyFolder}>
          Klasör İçeriği Seçin
        </Button>
        <label>{copyFolder ? copyFolder : "Klasör seçilmedi..."}</label>
      </div>
      <div className="flex flex-col gap-2 text-center">
        <Button variant={"outline"} onClick={handleTargetFolder}>
          Hedef Klasör Seçin
        </Button>
        <label>{targetFolder ? targetFolder : "Klasör seçilmedi..."}</label>
      </div>
      <div className="flex flex-col gap-2 text-center">
        <Button
          onClick={() => {
            if (excelPath && targetFolder) {
              CreateFolders(excelPath, copyFolder, targetFolder);
            }else{
                SendNotification("Excel Dosyası veya Hedef Klasör Seçilmedi", "", "", "warning");
            }
          }}
        >
          Onayla
        </Button>
      </div>
      
    </div>
  );
}

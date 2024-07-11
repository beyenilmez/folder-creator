import { useState } from "react";
import { Button } from "./ui/button";
import {
  CreateFoldersV2,
  GetCopyFolderDialog,
  GetExcelFileDialog,
  GetTargetFolderDialog,
  OpenFile,
  OpenFileInExplorer,
  SendNotification,
} from "@/wailsjs/go/main/App";
import { LogDebug } from "@/wailsjs/runtime/runtime";
import { LoaderCircle, X } from "lucide-react";

export function HomeV2() {
  const [excelPath, setExcelPath] = useState<string>("");
  const [copyFolder, setCopyFolder] = useState<string>("");
  const [targetFolder, setTargetFolder] = useState<string>("");

  const [running, setRunning] = useState<boolean>(false);
  const [message, setMessage] = useState<string>("");

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

  const handleRun = () => {
    setRunning(true);

    if (excelPath && targetFolder) {
      CreateFoldersV2(excelPath, copyFolder, targetFolder)
        .then((error) => {
          if (error !== "") {
            SendNotification("Hata", error, "", "error");
          }
        })
        .finally(() => {
          setRunning(false);
          setMessage("Klasörler başarıyla oluşturuldu");
        });
    } else {
      SendNotification(
        "Excel Dosyası veya Hedef Klasör Seçilmedi",
        "",
        "",
        "warning"
      );
    }
  };

  window.setExcelMessage = (message: string) => {
    LogDebug("window.setExcelMessage: " + message);

    setMessage(message);
  };

  window.setExcelProgress = (progress: number) => {
    LogDebug("window.setExcelProgress: " + progress);
  };

  return (
    <div className="flex flex-col justify-center items-center gap-5 w-full h-full">
      <div className="flex flex-col items-center gap-2 w-full">
        <Button variant={"outline"} onClick={handleExcelFileDialog}>
          Excel Dosyası Seçin
        </Button>
        <div className="flex h-4">
          <Button
            className="h-full"
            disabled={!excelPath}
            variant={"link"}
            onClick={() => OpenFile(excelPath)}
          >
            {excelPath ? excelPath : "Dosya seçilmedi..."}
          </Button>

          <Button
            variant={"destructive"}
            className={`rounded-sm w-4 h-4 ${!excelPath ? "hidden" : ""}`}
            size={"icon"}
            onClick={() => {
              setExcelPath("");
            }}
          >
            <X className="p-0.5" />
          </Button>
        </div>
      </div>

      <div className="flex flex-col items-center gap-2 w-full">
        <Button variant={"outline"} onClick={handleCopyFolder}>
          Klasör İçeriği Seçin
        </Button>
        <div className="flex h-4">
          <Button
            className="h-full"
            disabled={!copyFolder}
            variant={"link"}
            onClick={() => OpenFileInExplorer(copyFolder)}
          >
            {copyFolder ? copyFolder : "Klasör seçilmedi..."}
          </Button>
          <Button
            variant={"destructive"}
            className={`rounded-sm w-4 h-4 ${!copyFolder ? "hidden" : ""}`}
            size={"icon"}
            onClick={() => {
              setCopyFolder("");
            }}
          >
            <X className="p-0.5" />
          </Button>
        </div>
      </div>
      <div className="flex flex-col items-center gap-2 w-full">
        <Button variant={"outline"} onClick={handleTargetFolder}>
          Hedef Klasör Seçin
        </Button>
        <div className="flex h-4">
          <Button
            className="h-full"
            disabled={!targetFolder}
            variant={"link"}
            onClick={() => OpenFileInExplorer(targetFolder)}
          >
            {targetFolder ? targetFolder : "Klasör seçilmedi..."}
          </Button>
          <Button
            variant={"destructive"}
            className={`rounded-sm w-4 h-4 ${!targetFolder ? "hidden" : ""}`}
            size={"icon"}
            onClick={() => {
              setTargetFolder("");
            }}
          >
            <X className="p-0.5" />
          </Button>
        </div>
      </div>

      <div className="flex flex-col gap-2 text-center">
        <Button
          onClick={handleRun}
          className="mt-4 w-64"
          disabled={!excelPath || !targetFolder || running}
        >
          {running ? (
            <LoaderCircle className="w-6 h-6 animate-spin" />
          ) : (
            "Klasörleri Oluştur"
          )}
        </Button>
      </div>
      <div className="h-8 text-lg">{message}</div>
    </div>
  );
}

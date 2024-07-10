import { useEffect, useState } from "react";
import { Button } from "./ui/button";
import {
  CreateFolders,
  GetCopyFolderDialog,
  GetExcelFileDialog,
  GetTargetFolderDialog,
  GetWordFileDialog,
  SendNotification,
} from "@/wailsjs/go/main/App";
import { Input } from "./ui/input";
import { useConfig } from "@/contexts/config-provider";
import { LogDebug } from "@/wailsjs/runtime/runtime";
import { LoaderCircle } from "lucide-react";

export function Home() {
  const { config, setConfigField } = useConfig();

  const [excelPath, setExcelPath] = useState<string>("");
  const [wordPath, setWordPath] = useState<string>("");
  const [copyFolder, setCopyFolder] = useState<string>("");
  const [targetFolder, setTargetFolder] = useState<string>("");
  const [folderNamePattern, setFolderNamePattern] = useState<string>("");
  const [wordFileNamePattern, setWordFileNamePattern] = useState<string>("");

  const [running, setRunning] = useState<boolean>(false);
  const [message, setMessage] = useState<string>("");

  useEffect(() => {
    setFolderNamePattern(config?.folderNamePattern!);
    setWordFileNamePattern(config?.wordFileNamePattern!);
  }, [config]);

  const handleExcelFileDialog = () => {
    GetExcelFileDialog().then((path) => {
      setExcelPath(path);
    });
  };

  const handleWordFileDialog = () => {
    GetWordFileDialog().then((path) => {
      setWordPath(path);
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

    if (excelPath && targetFolder && folderNamePattern) {
      CreateFolders(
        excelPath,
        wordPath,
        copyFolder,
        targetFolder,
        folderNamePattern,
        wordFileNamePattern
      )
        .then((error) => {
          if (error !== "") {
            SendNotification("Hata", error, "", "error");
          }
        })
        .finally(() => {
          setRunning(false);
        });
    } else {
      SendNotification(
        "Excel Dosyası, Hedef Klasör veya Klasör Adı Seçilmedi",
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
        <label>{excelPath ? excelPath : "Dosya seçilmedi..."}</label>
      </div>
      <div className="flex flex-col items-center gap-2 w-full">
        <Button variant={"outline"} onClick={handleWordFileDialog}>
          Word Dosyası Seçin
        </Button>
        <label>{wordPath ? wordPath : "Dosya seçilmedi..."}</label>
      </div>
      <div className="flex flex-col items-center gap-2 w-full">
        <Button variant={"outline"} onClick={handleCopyFolder}>
          Klasör İçeriği Seçin
        </Button>
        <label>{copyFolder ? copyFolder : "Klasör seçilmedi..."}</label>
      </div>
      <div className="flex flex-col items-center gap-2 w-full">
        <Button variant={"outline"} onClick={handleTargetFolder}>
          Hedef Klasör Seçin
        </Button>
        <label>{targetFolder ? targetFolder : "Klasör seçilmedi..."}</label>
      </div>
      <div className="flex items-center gap-0 w-full">
        <div className="flex flex-col items-center gap-2 w-full">
          <label>Klasör Adı</label>
          <Input
            className="w-[90%]"
            value={folderNamePattern}
            onChange={(e) => {
              setConfigField("folderNamePattern", e.target.value);
              setFolderNamePattern(e.target.value);
            }}
          />
        </div>
        <div className="flex flex-col items-center gap-2 w-full">
          <label>Word Dosyası Adı</label>
          <Input
            className="w-[90%]"
            value={wordFileNamePattern}
            onChange={(e) => {
              setConfigField("wordFileNamePattern", e.target.value);
              setWordFileNamePattern(e.target.value);
            }}
          />
        </div>
      </div>
      <div className="flex flex-col gap-2 text-center">
        <Button
          onClick={handleRun}
          className="w-64"
          disabled={
            !excelPath || !targetFolder || !folderNamePattern || running
          }
        >
          {running ? (
            <LoaderCircle className="w-6 h-6 animate-spin" />
          ) : (
            "Onayla"
          )}
        </Button>
      </div>
      <div className="h-8 text-lg">{message}</div>
    </div>
  );
}

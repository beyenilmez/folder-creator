import { useEffect, useState } from "react";
import { Button } from "./ui/button";
import {
  CreateFolders,
  GetCopyFolderDialog,
  GetExcelFileDialog,
  GetFileDialog,
  GetTargetFolderDialog,
  GetWordFileDialog,
  OpenFile,
  OpenFileInExplorer,
  SendNotification,
} from "@/wailsjs/go/main/App";
import { Input } from "./ui/input";
import { useConfig } from "@/contexts/config-provider";
import { LogDebug } from "@/wailsjs/runtime/runtime";
import { LoaderCircle, X } from "lucide-react";

export function Home() {
  const { config, setConfigField } = useConfig();

  const [excelPath, setExcelPath] = useState<string>("");
  const [wordPath, setWordPath] = useState<string>("");
  const [filePath, setFilePath] = useState<string>("");
  const [copyFolder, setCopyFolder] = useState<string>("");
  const [targetFolder, setTargetFolder] = useState<string>("");
  const [folderNamePattern, setFolderNamePattern] = useState<string>("");
  const [wordFileNamePattern, setWordFileNamePattern] = useState<string>("");
  const [fileNamePattern, setFileNamePattern] = useState<string>("");

  const [running, setRunning] = useState<boolean>(false);
  const [message, setMessage] = useState<string>("");

  useEffect(() => {
    setFolderNamePattern(config?.folderNamePattern!);
    setWordFileNamePattern(config?.wordFileNamePattern!);
    setFileNamePattern(config?.fileNamePattern!);
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

  const handleFileDialog = () => {
    GetFileDialog().then((path) => {
      setFilePath(path);
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
        wordFileNamePattern,
        fileNamePattern,
        filePath
      )
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
        <Button variant={"outline"} onClick={handleWordFileDialog}>
          Word Dosyası Seçin
        </Button>
        <div className="flex h-4">
          <Button
            className="h-full"
            disabled={!wordPath}
            variant={"link"}
            onClick={() => OpenFile(wordPath)}
          >
            {wordPath ? wordPath : "Dosya seçilmedi..."}
          </Button>
          <Button
            variant={"destructive"}
            className={`rounded-sm w-4 h-4 ${!wordPath ? "hidden" : ""}`}
            size={"icon"}
            onClick={() => {
              setWordPath("");
            }}
          >
            <X className="p-0.5" />
          </Button>
        </div>
      </div>

      <div className="flex flex-col items-center gap-2 w-full">
        <Button variant={"outline"} onClick={handleFileDialog}>
          Bilirkişi Dilekçesi Seçin
        </Button>
        <div className="flex h-4">
          <Button
            className="h-full"
            disabled={!filePath}
            variant={"link"}
            onClick={() => OpenFile(filePath)}
          >
            {filePath ? filePath : "Dosya seçilmedi..."}
          </Button>
          <Button
            variant={"destructive"}
            className={`rounded-sm w-4 h-4 ${!filePath ? "hidden" : ""}`}
            size={"icon"}
            onClick={() => {
              setFilePath("");
            }}
          >
            <X className="p-0.5" />
          </Button>
        </div>
      </div>

      <div className="flex flex-col items-center gap-2 w-full">
        <Button variant={"outline"} onClick={handleCopyFolder}>
          Kopyalanacak Klasör Seçin
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
        <div className="flex flex-col items-center gap-2 w-full">
          <label>Dosya Adı</label>
          <Input
            className="w-[90%]"
            value={fileNamePattern}
            onChange={(e) => {
              setConfigField("fileNamePattern", e.target.value);
              setFileNamePattern(e.target.value);
            }}
          />
        </div>
      </div>
      <div className="flex flex-col gap-2 text-center">
        <Button
          onClick={handleRun}
          className="mt-4 w-64"
          disabled={
            !excelPath || !targetFolder || !folderNamePattern || running
          }
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

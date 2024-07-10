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

export function Home() {
  const { config, setConfigField } = useConfig();

  const [excelPath, setExcelPath] = useState<string>("");
  const [wordPath, setWordPath] = useState<string>("");
  const [copyFolder, setCopyFolder] = useState<string>("");
  const [targetFolder, setTargetFolder] = useState<string>("");
  const [folderNamePattern, setFolderNamePattern] = useState<string>("");
  const [wordFileNamePattern, setWordFileNamePattern] = useState<string>("");

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

  return (
    <div className="flex flex-col justify-center items-center gap-6 w-full h-full">
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
      <div className="flex flex-col items-center gap-2 w-full">
        <label>Klasör Adı</label>
        <Input
          className="w-full max-w-[30rem]"
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
          className="w-full max-w-[30rem]"
          value={wordFileNamePattern}
          onChange={(e) => {
            setConfigField("wordFileNamePattern", e.target.value);
            setWordFileNamePattern(e.target.value);
          }}
        />
      </div>
      <div className="flex flex-col gap-2 text-center">
        <Button
          onClick={() => {
            if (excelPath && targetFolder && folderNamePattern) {
              CreateFolders(
                excelPath,
                wordPath,
                copyFolder,
                targetFolder,
                folderNamePattern,
                wordFileNamePattern
              ).then((error) => {
                if (error !== "") {
                  SendNotification("Hata", error, "", "error");
                }
              });
            } else {
              SendNotification(
                "Excel Dosyası, Hedef Klasör veya Klasör Adı Seçilmedi",
                "",
                "",
                "warning"
              );
            }
          }}
        >
          Onayla
        </Button>
      </div>
    </div>
  );
}

import { useEffect, useState } from "react";
import { Button } from "./ui/button";
import {
  AddTapuToExcel,
  GetExcelFileDialog,
  GetTargetFolderDialog,
  OpenFile,
  OpenFileInExplorer,
} from "@/wailsjs/go/main/App";
import { LoaderCircle, X } from "lucide-react";
import { LogDebug } from "@/wailsjs/runtime/runtime";
import { Input } from "./ui/input";
import { useConfig } from "@/contexts/config-provider";

export function CiltSayfa() {
  const { config, setConfigField } = useConfig();

  const [ciltCellName, setCiltCellName] = useState<string>("");
  const [sayfaCellName, setSayfaCellName] = useState<string>("");
  const [mevkiCellName, setMevkiCellName] = useState<string>("");
  const [alanCellNameTapu, setAlanCellNameTapu] = useState<string>("");

  const [excelPath, setExcelPath] = useState<string>("");
  const [folderPath, setFolderPath] = useState<string>("");
  const [tapuNamePattern, setTapuNamePattern] = useState<string>("");

  const [message, setMessage] = useState<string>("");
  const [running, setRunning] = useState<boolean>(false);

  useEffect(() => {
    setCiltCellName(config?.ciltCellName!);
    setSayfaCellName(config?.sayfaCellName!);
    setTapuNamePattern(config?.tapuNamePattern!);
    setMevkiCellName(config?.mevkiCellName!);
    setAlanCellNameTapu(config?.alanCellNameTapu!);
  }, [config]);

  window.setCiltMessage = (message: string) => {
    LogDebug("window.setCiltMessage: " + message);

    setMessage(message);
  };

  const handleExcelFileDialog = () => {
    GetExcelFileDialog().then((path) => {
      setExcelPath(path);
    });
  };

  const handlefolderPath = () => {
    GetTargetFolderDialog().then((path) => {
      setFolderPath(path);
    });
  };

  const handleRun = () => {
    setRunning(true);
    AddTapuToExcel(
      excelPath,
      folderPath,
      tapuNamePattern,
      ciltCellName,
      sayfaCellName,
      mevkiCellName,
      alanCellNameTapu
    )
      .then((error) => {
        if (error !== "") {
          setMessage(error);
        }
      })
      .finally(() => {
        setRunning(false);
      });
  };

  return (
    <div className="flex flex-col justify-center items-center gap-12 w-full h-full">
      <div className="flex flex-row">
        <div className="flex flex-col items-center gap-2 w-full">
          <label>Cilt Sütunu</label>
          <Input
            className="w-[90%]"
            value={ciltCellName}
            onChange={(e) => {
              setConfigField("ciltCellName", e.target.value);
              setCiltCellName(e.target.value);
            }}
          />
        </div>

        <div className="flex flex-col items-center gap-2 w-full">
          <label>Sayfa Sütunu</label>
          <Input
            className="w-[90%]"
            value={sayfaCellName}
            onChange={(e) => {
              setConfigField("sayfaCellName", e.target.value);
              setSayfaCellName(e.target.value);
            }}
          />
        </div>

        <div className="flex flex-col items-center gap-2 w-full">
          <label>Mevki Sütunu</label>
          <Input
            className="w-[90%]"
            value={mevkiCellName}
            onChange={(e) => {
              setConfigField("mevkiCellName", e.target.value);
              setMevkiCellName(e.target.value);
            }}
          />
        </div>

        <div className="flex flex-col items-center gap-2 w-full">
          <label>Alan Sütunu</label>
          <Input
            className="w-[90%]"
            value={alanCellNameTapu}
            onChange={(e) => {
              setConfigField("alanCellNameTapu", e.target.value);
              setAlanCellNameTapu(e.target.value);
            }}
          />
        </div>
      </div>

      <div className="flex flex-col items-center gap-2 w-full">
        <Button variant={"outline"} onClick={handlefolderPath}>
          Klasör Seçin
        </Button>
        <div className="flex h-4">
          <Button
            className="h-full"
            disabled={!folderPath}
            variant={"link"}
            onClick={() => OpenFileInExplorer(folderPath)}
          >
            {folderPath ? folderPath : "Klasör seçilmedi..."}
          </Button>
          <Button
            variant={"destructive"}
            className={`rounded-sm w-4 h-4 ${!folderPath ? "hidden" : ""}`}
            size={"icon"}
            onClick={() => {
              setFolderPath("");
            }}
          >
            <X className="p-0.5" />
          </Button>
        </div>
      </div>

      <div className="flex flex-col items-center gap-2 w-full">
        <label>Evrak dosyası konumunu seçin</label>
        <Input
          className="w-1/2"
          value={tapuNamePattern}
          onChange={(e) => {
            setConfigField("tapuNamePattern", e.target.value);
            setTapuNamePattern(e.target.value);
          }}
        />
      </div>

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

      <div className="flex flex-col gap-2 text-center">
        <Button
          disabled={
            !folderPath ||
            !tapuNamePattern ||
            !(
              sayfaCellName ||
              ciltCellName ||
              mevkiCellName ||
              alanCellNameTapu
            ) ||
            !excelPath ||
            running
          }
          onClick={handleRun}
          className="w-64"
        >
          {running ? (
            <LoaderCircle className="w-6 h-6 animate-spin" />
          ) : (
            "Sütunları Ekle"
          )}
        </Button>
      </div>
      <div className="h-8 text-lg">{message}</div>
    </div>
  );
}

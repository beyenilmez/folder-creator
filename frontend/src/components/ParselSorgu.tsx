import { useEffect, useState } from "react";
import { Button } from "./ui/button";
import {
  AddParselSorguFields,
  GetExcelFileDialog,
  OpenFile,
} from "@/wailsjs/go/main/App";
import { LoaderCircle, X } from "lucide-react";
import { LogDebug } from "@/wailsjs/runtime/runtime";
import { Input } from "./ui/input";
import { useConfig } from "@/contexts/config-provider";
import { Switch } from "./ui/switch";

export function ParselSorguComp() {
  const { config, setConfigField } = useConfig();

  const [ilCellName, setIlCellName] = useState<string>("");
  const [ilceCellName, setIlceCellName] = useState<string>("");
  const [mahalleCellName, setMahalleCellName] = useState<string>("");
  const [adaCellName, setAdaCellName] = useState<string>("");
  const [parselCellName, setParselCellName] = useState<string>("");
  const [alanCellName, setAlanCellName] = useState<string>("");
  const [paftaCellName, setPaftaCellName] = useState<string>("");
  const [parselSorguHeadless, setParselSorguHeadless] =
    useState<boolean>(false);

  const [excelPath, setExcelPath] = useState<string>("");

  const [message, setMessage] = useState<string>("");
  const [running, setRunning] = useState<boolean>(false);

  useEffect(() => {
    setIlCellName(config?.ilCellName!);
    setIlceCellName(config?.ilceCellName!);
    setMahalleCellName(config?.mahalleCellName!);
    setAdaCellName(config?.adaCellName!);
    setParselCellName(config?.parselCellName!);
    setAlanCellName(config?.alanCellName!);
    setPaftaCellName(config?.paftaCellName!);
    setParselSorguHeadless(config?.parselSorguHeadless!);
  }, [config]);

  window.setParselMessage = (message: string) => {
    LogDebug("window.setParselMessage: " + message);

    setMessage(message);
  };

  const handleExcelFileDialog = () => {
    GetExcelFileDialog().then((path) => {
      setExcelPath(path);
    });
  };

  const handleHeadless = () => {
    setConfigField("parselSorguHeadless", !parselSorguHeadless);
    setParselSorguHeadless(!parselSorguHeadless);
  };

  const handleRun = () => {
    setRunning(true);
    AddParselSorguFields(
      excelPath,
      ilCellName,
      ilceCellName,
      mahalleCellName,
      adaCellName,
      parselCellName,
      alanCellName,
      paftaCellName,
      parselSorguHeadless
    ).finally(() => {
      setRunning(false);
    });
  };

  return (
    <div className="flex flex-col justify-center items-center gap-12 w-full h-full">
      <div className="flex flex-row">
        <div className="flex flex-col items-center gap-2 w-full">
          <label>İl Sütunu</label>
          <Input
            className="w-[90%]"
            value={ilCellName}
            onChange={(e) => {
              setConfigField("ilCellName", e.target.value);
              setIlCellName(e.target.value);
            }}
          />
        </div>

        <div className="flex flex-col items-center gap-2 w-full">
          <label>İlçe Sütunu</label>
          <Input
            className="w-[90%]"
            value={ilceCellName}
            onChange={(e) => {
              setConfigField("ilceCellName", e.target.value);
              setIlceCellName(e.target.value);
            }}
          />
        </div>

        <div className="flex flex-col items-center gap-2 w-full">
          <label>Mahalle Sütunu</label>
          <Input
            className="w-[90%]"
            value={mahalleCellName}
            onChange={(e) => {
              setConfigField("mahalleCellName", e.target.value);
              setMahalleCellName(e.target.value);
            }}
          />
        </div>

        <div className="flex flex-col items-center gap-2 w-full">
          <label>Ada Sütunu</label>
          <Input
            className="w-[90%]"
            value={adaCellName}
            onChange={(e) => {
              setConfigField("adaCellName", e.target.value);
              setAdaCellName(e.target.value);
            }}
          />
        </div>

        <div className="flex flex-col items-center gap-2 w-full">
          <label>Parsel Sütunu</label>
          <Input
            className="w-[90%]"
            value={parselCellName}
            onChange={(e) => {
              setConfigField("parselCellName", e.target.value);
              setParselCellName(e.target.value);
            }}
          />
        </div>

        <div className="flex flex-col items-center gap-2 w-full">
          <label>Alan Sütunu</label>
          <Input
            className="w-[90%]"
            value={alanCellName}
            onChange={(e) => {
              setConfigField("alanCellName", e.target.value);
              setAlanCellName(e.target.value);
            }}
          />
        </div>

        <div className="flex flex-col items-center gap-2 w-full">
          <label>Pafta Sütunu</label>
          <Input
            className="w-[90%]"
            value={paftaCellName}
            onChange={(e) => {
              setConfigField("paftaCellName", e.target.value);
              setPaftaCellName(e.target.value);
            }}
          />
        </div>
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
        <div className="flex flex-col items-center gap-8">
          <div className="flex flex-col items-center gap-2 font-medium text-lg">
            Gizli
            <Switch
              checked={parselSorguHeadless}
              onCheckedChange={handleHeadless}
            />
          </div>
          <Button
            disabled={
              !ilCellName ||
              !ilceCellName ||
              !mahalleCellName ||
              !adaCellName ||
              !parselCellName ||
              !(alanCellName || paftaCellName) ||
              !excelPath ||
              running
            }
            onClick={handleRun}
            className="w-64"
          >
            {running ? (
              <LoaderCircle className="w-6 h-6 animate-spin" />
            ) : (
              "Alan ve Pafta Ekle"
            )}
          </Button>
        </div>
        <div className="h-8 text-lg">{message}</div>
      </div>
    </div>
  );
}

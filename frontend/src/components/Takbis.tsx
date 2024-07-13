import { useEffect, useState } from "react";
import { Button } from "./ui/button";
import {
  GetExcelFileDialog,
  GetExcelFilesDialog,
  ModifyExcelWithTakbis,
  OpenFile,
} from "@/wailsjs/go/main/App";
import { LoaderCircle, X } from "lucide-react";
import { useConfig } from "@/contexts/config-provider";
import { Input } from "./ui/input";
import { LogDebug } from "@/wailsjs/runtime/runtime";

export function Takbis() {
  const { config, setConfigField } = useConfig();

  const [takbisPaths, setTakbisPaths] = useState<string[]>([]);
  const [excelPath, setExcelPath] = useState<string>("");
  const [excelHeaderMatchPattern, setExcelHeaderMatchPattern] =
    useState<string>("");
  const [excelCellModifyPattern, setExcelCellModifyPattern] =
    useState<string>("");

  const [message, setMessage] = useState<string>("");
  const [running, setRunning] = useState<boolean>(false);

  useEffect(() => {
    setExcelHeaderMatchPattern(config?.excelHeaderMatchPattern!);
    setExcelCellModifyPattern(config?.excelCellModifyPattern!);
  }, [config]);

  const handleExcelFilesDialog = () => {
    GetExcelFilesDialog().then((paths) => {
      if (paths) {
        setTakbisPaths(paths);
      } else {
        setTakbisPaths([]);
      }
    });
  };

  const handleExcelFileDialog = () => {
    GetExcelFileDialog().then((path) => {
      setExcelPath(path);
    });
  };

  const handleRun = () => {
    setRunning(true);
    ModifyExcelWithTakbis(
      excelPath,
      takbisPaths,
      excelHeaderMatchPattern,
      excelCellModifyPattern
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

  window.setTakbisMessage = (message: string) => {
    LogDebug("window.setTakbisMessage: " + message);

    setMessage(message);
  };

  return (
    <div className="flex flex-col gap-12 p-16 w-full h-full">
      <div className="flex gap-64">
        <div className="flex flex-col items-center gap-2 w-full">
          <Button variant={"outline"} onClick={handleExcelFilesDialog}>
            Takbis Dosyalarını Seçin
          </Button>
          <div className="flex flex-col">
            {takbisPaths.length > 0 ? (
              takbisPaths.map((path) => (
                <div className="flex items-center h-3.5" key={path}>
                  <Button className="h-full" variant={"link"} onClick={() => OpenFile(path)}>
                    {" "}
                    {path}
                  </Button>

                  <Button
                    variant={"destructive"}
                    className={`rounded-sm w-4 h-4`}
                    size={"icon"}
                    onClick={() => {
                      setTakbisPaths(takbisPaths.filter((p) => p !== path));
                    }}
                  >
                    <X className="p-0.5" />
                  </Button>
                </div>
              ))
            ) : (
              <Button
                className="flex items-center h-5"
                variant={"link"}
                disabled
              >
                Dosya seçilmedi...
              </Button>
            )}
          </div>
        </div>

        <div className="flex flex-col items-center gap-2 w-full">
          <Button variant={"outline"} onClick={handleExcelFileDialog}>
            Hedef Excel Dosyası Seçin
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
      </div>

      <div className="flex flex-col items-center gap-2 w-full">
        <label>Sütun Adı Karşılaştırma Kuralı</label>
        <Input
          className="w-[90%]"
          value={excelHeaderMatchPattern}
          onChange={(e) => {
            setConfigField("excelHeaderMatchPattern", e.target.value);
            setExcelHeaderMatchPattern(e.target.value);
          }}
        />
      </div>

      <div className="flex flex-col items-center gap-2 w-full">
        <label>Hücre Değişimi Kuralı</label>
        <Input
          className="w-[90%]"
          value={excelCellModifyPattern}
          onChange={(e) => {
            setConfigField("excelCellModifyPattern", e.target.value);
            setExcelCellModifyPattern(e.target.value);
          }}
        />
      </div>

      <div className="flex flex-col items-center gap-2 w-full">
        <Button
          disabled={
            !excelCellModifyPattern ||
            !excelHeaderMatchPattern ||
            takbisPaths.length === 0 ||
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
        <div className="h-8 text-lg">{message}</div>
      </div>
    </div>
  );
}

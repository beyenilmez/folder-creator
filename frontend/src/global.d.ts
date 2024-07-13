declare global {
  interface Window {
    toast: (props: ToastProps) => void;
    goto: goto;
    setExcelMessage: (message: string) => void;
    setParselMessage: (message: string) => void;
    setExcelProgress: (progress: number) => void;
    setCiltMessage: (message: string) => void;
  }
}

export {};

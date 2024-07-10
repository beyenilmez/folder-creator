declare global {
  interface Window {
    toast: (props: ToastProps) => void;
    goto: goto;
    setExcelMessage: (message: string) => void;
    setExcelProgress: (progress: number) => void;
  }
}

export {};

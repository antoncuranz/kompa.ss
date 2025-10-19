import React, { createContext, useContext, ReactNode } from 'react';
import * as DialogPrimitive from "@radix-ui/react-dialog"
import {useRouter} from "next/navigation";
import {DialogContent} from "@/components/ui/dialog.tsx"

interface DialogContextType {
  onClose: (needsUpdate?: boolean) => void;
}

const DialogContext = createContext<DialogContextType | undefined>(undefined);

export function DialogContextProvider({
  children, onClose
}: {
  onClose: (needsUpdate?: boolean) => void
  children: ReactNode
}) {

  return (
      <DialogContext.Provider value={{ onClose }}>
        {children}
      </DialogContext.Provider>
  );
}

export function useDialogContext(): DialogContextType {
  const context = useContext(DialogContext);
  if (!context) {
    throw new Error('useUser must be used within a UserProvider');
  }
  return context;
}

export function Dialog({
  children, open, setOpen
}: {
  children: React.ReactNode | React.ReactNode[]
  open: boolean
  setOpen: (needsUpdate: boolean) => void
}) {
  const router = useRouter()

  function onClose(needsUpdate?: boolean) {
    setOpen(false)
    if (needsUpdate)
      router.refresh()
  }

  return (
      <DialogPrimitive.Root open={open} onOpenChange={open => open || onClose(false)}>
        <DialogContent>
          <DialogContextProvider onClose={onClose}>
            {children}
          </DialogContextProvider>
        </DialogContent>
      </DialogPrimitive.Root>
  )
}
Dialog.displayName = DialogPrimitive.Root.displayName

export const RowContainer = ({
  children
}: {
  children: React.ReactNode;
}) => {
  return (
      <div className="grid md:grid-cols-2 gap-2">
        {children}
      </div>
  );
};

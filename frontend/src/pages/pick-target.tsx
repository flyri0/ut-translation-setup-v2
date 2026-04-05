import { useEffect, useState } from "react"
import { FaFileArrowDown, FaFolderOpen, FaFile } from "react-icons/fa6"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import {
  Field,
  FieldDescription,
  FieldLabel,
  FieldGroup,
  FieldContent,
  FieldTitle,
} from "@/components/ui/field"
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group"
import { ButtonGroup } from "@/components/ui/button-group"
import {
  QuickFind,
  OpenFilePicker,
  CheckFreeSpace,
  SaveSettings,
} from "../../wailsjs/go/main/PickTargetService"
import { toast } from "sonner"
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { Toaster } from "@/components/ui/sonner"
import { useLocation } from "wouter"

export default function PickTargetPage() {
  const [, navigate] = useLocation()
  const [targetPath, setTargetPath] = useState<string>("")
  const [isDemo, setIsDemo] = useState<boolean>(false)
  const [isBackupChecked, setIsBackupChecked] = useState<boolean>(false)
  const [isSpaceDialogOpen, setIsSpaceDialogOpen] = useState<boolean>(false)
  const [isValid, setIsValid] = useState<boolean>(false)
  const [isCheckingSpace, setIsCheckingSpace] = useState<boolean>(false)

  async function handleInstallButton() {
    SaveSettings(targetPath, isDemo, isBackupChecked)
    navigate("/install")
  }

  async function handleSuccessfulValidation(
    path: string,
    autoFound: boolean,
    isDemo: boolean
  ) {
    setTargetPath(path)
    setIsDemo(isDemo)

    const hasSpace = await CheckFreeSpace(path)

    if (!hasSpace) {
      setIsSpaceDialogOpen(true)
      setIsValid(false)
    } else {
      setIsValid(true)
      if (autoFound) {
        toast.success("Instalação encontrada automaticamente", {
          description: "Detectamos uma instalação do Until Then na Steam.",
          position: "top-center",
        })
      }
    }
  }

  async function handleManualSelect() {
    try {
      const result = await OpenFilePicker()
      if (result.valid) {
        handleSuccessfulValidation(result.path, false, result.isDemo)
      } else if (result.error !== "not_selected") {
        toast.error("Arquivo inválido", {
          position: "top-center",
        })
      }
    } catch {
      // TODO: colocar log aqui
    }
  }

  async function handleRetrySpaceCheck() {
    setIsCheckingSpace(true)
    const hasSpace = await CheckFreeSpace(targetPath)
    setIsCheckingSpace(false)

    if (hasSpace) {
      setIsSpaceDialogOpen(false)
      setIsValid(true)
    } else {
      toast.error("Ainda não há espaço suficiente disponivel.", {
        position: "top-center",
      })
    }
  }

  const handleCancelSpaceCheck = () => {
    setIsSpaceDialogOpen(false)
    setTargetPath("")
    setIsValid(false)
  }

  useEffect(() => {
    const runQuickFind = async () => {
      try {
        const result = await QuickFind()
        if (result.valid) {
          handleSuccessfulValidation(result.path, true, result.isDemo)
        }
      } catch {
        // TODO: put log here
      }
    }

    runQuickFind()
  }, [])

  return (
    <>
      <Toaster />
      <div className="mx-17 flex h-dvh flex-col justify-center select-none">
        <FieldGroup>
          <Field>
            <FieldLabel>Instalação do jogo</FieldLabel>
            <FieldDescription>
              Selecione o arquivo UntilThen.pck na pasta do jogo
            </FieldDescription>
            <ButtonGroup>
              <InputGroup>
                <InputGroupAddon
                  align="inline-start"
                  className="pointer-events-none"
                >
                  <FaFile />
                </InputGroupAddon>
                <InputGroupInput
                  disabled
                  value={targetPath}
                  placeholder="Nenhum arquivo selecionado"
                  className="truncate"
                />
              </InputGroup>

              <Button
                className="hover:cursor-pointer"
                variant="outline"
                onClick={handleManualSelect}
              >
                <FaFolderOpen />
                Procurar...
              </Button>
            </ButtonGroup>

            <FieldDescription className="h-5">
              {isValid && (
                <span className="font-medium text-emerald-500">
                  Arquivo UntilThen.pck validado com sucesso.
                </span>
              )}
            </FieldDescription>
          </Field>

          <Field orientation="horizontal">
            <Checkbox
              checked={isBackupChecked}
              onCheckedChange={(value: boolean) => setIsBackupChecked(value)}
            />
            <FieldContent>
              <FieldTitle>Criar backup de segurança</FieldTitle>
              <FieldDescription className="w-3/5">
                Salva uma cópia do arquivo original para restaurar o idioma
                padrão quando necessário. (Cerca de 2GB a mais vão ser ocupados)
              </FieldDescription>
            </FieldContent>
          </Field>
        </FieldGroup>

        <Button
          disabled={!isValid}
          onClick={handleInstallButton}
          size="lg"
          className="mt-5 w-full hover:cursor-pointer"
        >
          <FaFileArrowDown />
          Iniciar Instalação
        </Button>

        <AlertDialog
          open={isSpaceDialogOpen}
          onOpenChange={setIsSpaceDialogOpen}
        >
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Espaço em disco insuficiente</AlertDialogTitle>
              <AlertDialogDescription>
                O local selecionado não tem os <strong>5GB</strong> de espaço
                livre necessários. Libere espaço no disco e tente novamente.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <Button variant="outline" onClick={handleCancelSpaceCheck}>
                Cancelar
              </Button>
              <Button
                disabled={isCheckingSpace}
                onClick={handleRetrySpaceCheck}
              >
                {isCheckingSpace ? "Verificando..." : "Tentar Novamente"}
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>
    </>
  )
}

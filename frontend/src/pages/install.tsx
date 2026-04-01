import { useEffect, useRef, useState } from "react"
import { FaCircleCheck, FaCircleXmark } from "react-icons/fa6"
import { useLocation } from "wouter"

import { Item, ItemContent, ItemMedia, ItemTitle } from "@/components/ui/item"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Spinner } from "@/components/ui/spinner"

import { EventsOn, EventsOff, Quit } from "../../wailsjs/runtime"
import { RunInstallation } from "../../wailsjs/go/main/PckExplorerService"

export default function InstallPage() {
  const MAX_LOGS = 25
  const [, navigate] = useLocation()
  const [step, setStep] = useState("Preparando instalação...")
  const [progress, setProgress] = useState(0)
  const [logs, setLogs] = useState<string[]>([])
  const [status, setStatus] = useState<"installing" | "success" | "error">(
    "installing"
  )
  const scrollViewportRef = useRef<HTMLDivElement>(null)

  const appendLog = (message: string) => {
    setLogs((prev) => {
      const newLogs = [...prev, message]
      return newLogs.length > MAX_LOGS ? newLogs.slice(-MAX_LOGS) : newLogs
    })
  }

  useEffect(() => {
    RunInstallation().catch((err) => {
      setStatus("error")
      setStep("Falha crítica ao iniciar.")
      appendLog(`[SISTEMA] Erro ao invocar backend: ${err}`)
    })

    EventsOn("install_step", (currentStep: string) => {
      setStep(currentStep)
    })

    EventsOn("unzip_bin_progress", (p: number) => setProgress(p))
    EventsOn("unzip_trans_progress", (p: number) => setProgress(p))

    EventsOn("install_log", (msg: string) => {
      appendLog(msg)
    })

    EventsOn("install_error", (err: string) => {
      setStatus("error")
      setStep("Erro durante a instalação.")
      appendLog(`[ERRO] ${err}`)
    })

    EventsOn("install_success", (msg: string) => {
      setStatus("success")
      setStep(msg)
      setProgress(100)
      appendLog("[SISTEMA] Processo finalizado com sucesso!")
    })

    return () => {
      EventsOff("install_step")
      EventsOff("unzip_bin_progress")
      EventsOff("unzip_trans_progress")
      EventsOff("install_log")
      EventsOff("install_error")
      EventsOff("install_success")
    }
  }, [])

  useEffect(() => {
    if (scrollViewportRef.current) {
      scrollViewportRef.current.scrollTop =
        scrollViewportRef.current.scrollHeight
    }
  }, [logs])

  return (
    <div className="mx-20 flex h-screen flex-col items-center justify-center gap-5">
      <Item
        className="w-1/2 rounded transition-all duration-300"
        variant="muted"
      >
        <ItemMedia>
          {status === "installing" && <Spinner />}
          {status === "success" && <FaCircleCheck className="text-green-500" />}
          {status === "error" && <FaCircleXmark className="text-red-500" />}
        </ItemMedia>
        <ItemContent>
          <ItemTitle className={status === "error" ? "text-red-500" : ""}>
            {step}
          </ItemTitle>
        </ItemContent>
        <ItemContent className="flex-none justify-end font-mono text-sm">
          {progress}%
        </ItemContent>
      </Item>

      <ScrollArea
        ref={scrollViewportRef}
        className="h-1/3 w-full rounded border bg-black/5 p-3 font-mono text-xs text-muted-foreground dark:bg-black/40"
      >
        <div className="flex flex-col gap-1">
          {logs.length === 0 ? (
            <span className="italic opacity-50">Aguardando logs...</span>
          ) : (
            logs.map((log, index) => (
              <div
                key={index}
                className={log.startsWith("[ERRO]") ? "text-red-500" : ""}
              >
                {log}
              </div>
            ))
          )}
        </div>
      </ScrollArea>

      {status !== "installing" && (
        <button
          className="mt-4 rounded bg-primary px-6 py-2 text-primary-foreground transition-opacity hover:opacity-90"
          onClick={() => {
            if (status === "success") {
              navigate("/finished")
              return
            }

            Quit()
          }}
        >
          {status === "success" ? "Concluir" : "Fechar Instalador"}
        </button>
      )}
    </div>
  )
}

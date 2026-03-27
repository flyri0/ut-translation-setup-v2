import { FaMagnifyingGlass, FaFolderOpen, FaFile } from "react-icons/fa6"
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
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip"
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group"

export default function PickFilePage() {
  return (
    <div className="mx-17 flex h-dvh flex-col justify-center select-none">
      <FieldGroup>
        <Field>
          <FieldLabel>Selecione o local de instalação do jogo</FieldLabel>
          <InputGroup className="pointer-events-none cursor-default">
            <InputGroupInput
              disabled
              placeholder="Nenhum diretório selecionado..."
            />
            <InputGroupAddon align="inline-start">
              <FaFile />
            </InputGroupAddon>
          </InputGroup>
          <FieldDescription>Nenhum caminho selecionado</FieldDescription>
        </Field>

        <Field orientation="horizontal">
          <Checkbox />
          <FieldContent>
            <FieldTitle>Fazer backup</FieldTitle>
            <FieldDescription className="w-3/5">
              Cria uma cópia dos arquivos originais (requer 4GB livres). Útil
              para restaurar o idioma padrão posteriormente
            </FieldDescription>
          </FieldContent>
        </Field>
      </FieldGroup>

      <div className="mt-5 flex w-full gap-2">
        <Tooltip>
          <TooltipTrigger asChild>
            <Button size="lg" className="flex-2 hover:cursor-pointer">
              <FaMagnifyingGlass />
              Busca Automática
            </Button>
          </TooltipTrigger>
          <TooltipContent side="bottom">
            Localiza o jogo automaticamente (suporta apenas versões da Steam)
          </TooltipContent>
        </Tooltip>

        <Tooltip>
          <TooltipTrigger asChild>
            <Button
              variant="secondary"
              size="lg"
              className="flex-1 hover:cursor-pointer"
            >
              <FaFolderOpen />
              Selecionar Manualmente
            </Button>
          </TooltipTrigger>
          <TooltipContent side="bottom">
            Selecione a pasta onde o executável do jogo está instalado
          </TooltipContent>
        </Tooltip>
      </div>
    </div>
  )
}

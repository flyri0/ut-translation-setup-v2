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

export default function PickTargetPage() {
  return (
    <div className="mx-17 flex h-dvh flex-col justify-center select-none">
      <FieldGroup>
        <Field>
          <FieldLabel>Local de instalação do jogo</FieldLabel>
          <FieldDescription>
            Selecione a pasta onde o arquivo executável do jogo está localizado
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
                placeholder="Nenhum local selecionado..."
              />
            </InputGroup>

            <Button className="hover:cursor-pointer" variant="outline">
              <FaFolderOpen />
              Procurar...
            </Button>
          </ButtonGroup>
          <FieldDescription />
        </Field>

        <Field orientation="horizontal">
          <Checkbox />
          <FieldContent>
            <FieldTitle>Criar backup de segurança</FieldTitle>
            <FieldDescription className="w-3/5">
              Recomendado para restaurar os arquivos originais e o idioma padrão
              caso necessário
            </FieldDescription>
          </FieldContent>
        </Field>
      </FieldGroup>

      <Button disabled size="lg" className="mt-5 w-full hover:cursor-pointer">
        <FaFileArrowDown />
        Iniciar Instalação
      </Button>
    </div>
  )
}

import { FaSun, FaMoon } from "react-icons/fa6"
import { Button } from "@/components/ui/button"
import {
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from "@/components/ui/tooltip"
import { cn } from "@/lib/utils"
import { useTheme } from "./theme-provider"

interface ThemeToggleProps {
  className?: string
}

export default function ThemeToggle({ className }: ThemeToggleProps) {
  const { setTheme, theme } = useTheme()

  function onClick() {
    if (theme == "light") {
      setTheme("dark")
    } else {
      setTheme("light")
    }
  }

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <Button
          variant="outline"
          onClick={onClick}
          className={cn("hover:cursor-pointer", className)}
          size="icon"
        >
          {theme == "light" ? <FaSun /> : <FaMoon />}
        </Button>
      </TooltipTrigger>
      <TooltipContent>Mudar tema</TooltipContent>
    </Tooltip>
  )
}

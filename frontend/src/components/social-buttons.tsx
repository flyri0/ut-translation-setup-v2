import { BrowserOpenURL } from "../../wailsjs/runtime"
import { FaDiscord, FaGithub } from "react-icons/fa6"
import { Button } from "./ui/button"
import { Tooltip, TooltipContent, TooltipTrigger } from "./ui/tooltip"
import { cn } from "@/lib/utils"

interface SocialButtonsProps {
  className?: string
}

export default function SocialButtons({ className }: SocialButtonsProps) {
  return (
    <div className={cn("flex w-full gap-1", className)}>
      <Tooltip disableHoverableContent>
        <TooltipTrigger asChild>
          <Button
            className="flex-1 hover:cursor-pointer"
            variant="outline"
            onClick={() => BrowserOpenURL("https://discord.gg/MKn6QBVG9g")}
          >
            <FaDiscord />
            Discord
          </Button>
        </TooltipTrigger>
        <TooltipContent side="bottom">Entrar na comunidade</TooltipContent>
      </Tooltip>

      <Tooltip disableHoverableContent>
        <TooltipTrigger asChild>
          <Button
            className="flex-1 hover:cursor-pointer"
            variant="outline"
            onClick={() =>
              BrowserOpenURL(
                "https://github.com/flyri0/ut-translation-setup-v2"
              )
            }
          >
            <FaGithub />
            Github
          </Button>
        </TooltipTrigger>
        <TooltipContent side="bottom">Ver código-fonte</TooltipContent>
      </Tooltip>
    </div>
  )
}

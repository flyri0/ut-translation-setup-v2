import { Button } from "@/components/ui/button"
import { FaUsers, FaDiscord, FaGithub } from "react-icons/fa6"
import { Link } from "wouter"
import { BrowserOpenURL } from "../../wailsjs/runtime"
import ThemeToggle from "@/components/theme-toggle"

export default function WelcomePage() {
  const TEAM_MEMBERS = [
    { alias: "PitterG4", name: "Bernardo Hoffmann" },
    { alias: "Percival", name: "Gabriel Araújo" },
    { alias: "Lucasxt", name: "Lucas Silva" },
    { alias: "Yubi", name: "Eduarda Albuquerque" },
    { alias: "Ceci", name: "Cecília" },
    { alias: "flyri0", name: "Francisco" },
  ]

  return (
    <div className="flex text-foreground select-none">
      <img className="h-dvh max-w-fit" draggable={false} src="banner.webp" />

      <div className="mx-10 my-10 flex grow flex-col justify-center">
        <ThemeToggle className="mb-5" />

        <h1 className="text-3xl font-medium tracking-tight text-balance">
          Until Then... <span className="italic">em português!</span> ✨
        </h1>

        <p className="mt-5 font-serif text-muted-foreground">
          Esta tradução foi feita com carinho por fãs, para que mais pessoas
          possam desfrutar de <span className="italic">Until Then</span> em
          nosso belíssimo idioma. Esperamos que esta jornada emocione você tanto
          quanto nos emocionou.
        </p>

        <div className="mt-10">
          <p className="mb-2 flex items-center gap-2 text-lg font-medium text-muted-foreground">
            <span>
              <FaUsers />
            </span>
            Equipe do Projeto:
          </p>
          <ul className="list-inside list-disc font-serif">
            {TEAM_MEMBERS.map(({ alias, name }) => (
              <li key={alias}>
                <span className="text-muted-foreground">({alias})</span> {name}
              </li>
            ))}
          </ul>
        </div>

        <Link to="/pick-target" asChild>
          <Button className="mt-10 w-full self-center border-2 py-10 text-4xl transition hover:cursor-pointer hover:border-2">
            Começar
          </Button>
        </Link>

        <div className="mt-2 flex w-full justify-center gap-2">
          <Button
            variant="outline"
            onClick={() => BrowserOpenURL("https://discord.gg/MKn6QBVG9g")}
            className="flex h-auto flex-1 flex-col items-center px-6 py-4 hover:cursor-pointer"
          >
            <div className="flex items-center gap-2 text-lg">
              <FaDiscord />
              Discord
            </div>
            <span className="text-[10px] uppercase opacity-80">
              Entre na comunidade
            </span>
          </Button>

          <Button
            variant="outline"
            onClick={() =>
              BrowserOpenURL(
                "https://github.com/flyri0/ut-translation-setup-v2"
              )
            }
            className="flex h-auto flex-1 flex-col items-center px-6 py-4 hover:cursor-pointer"
          >
            <div className="flex items-center gap-2 text-lg">
              <FaGithub />
              GitHub
            </div>
            <span className="text-[10px] uppercase opacity-80">
              Código-fonte
            </span>
          </Button>
        </div>
      </div>
    </div>
  )
}

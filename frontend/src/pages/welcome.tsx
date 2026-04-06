import { Button } from "@/components/ui/button"
import { FaUsers } from "react-icons/fa6"
import { Link } from "wouter"
import ThemeToggle from "@/components/theme-toggle"
import SocialButtons from "@/components/social-buttons"

export default function WelcomePage() {
  const TEAM_MEMBERS = [
    { alias: "PitterG4", name: "Bernardo Hoffmann" },
    { alias: "Percival", name: "Gabriel Araújo" },
    { alias: "Lucasxt", name: "Lucas Silva" },
    { alias: "Yubi", name: "Eduarda Albuquerque" },
    { alias: "Ceci", name: "Cecília" },
    { alias: "flyri0", name: "Francisco Lyrio" },
  ]

  return (
    <div className="flex">
      <img className="h-dvh" draggable={false} src="banner.webp" />

      <div className="mx-5 my-2 self-center">
        <ThemeToggle className="mb-1" />

        <h1 className="text-base font-medium tracking-tight text-balance md:text-2xl">
          Until Then... <span className="italic">em português!</span> ✨
        </h1>

        <p className="font-serif text-xs text-muted-foreground md:text-base">
          Esta tradução foi feita com carinho por fãs, para que mais pessoas
          possam desfrutar de <span className="italic">Until Then</span> em
          nosso belíssimo idioma. Esperamos que esta jornada emocione você tanto
          quanto nos emocionou.
        </p>

        <p className="mt-2 flex items-center gap-1 text-sm text-muted-foreground md:text-base">
          <span>
            <FaUsers />
          </span>
          Equipe do Projeto:
        </p>
        <ul className="list-inside list-disc font-serif text-xs md:text-sm">
          {TEAM_MEMBERS.map(({ alias, name }) => (
            <li key={alias}>
              <span className="text-muted-foreground">({alias})</span> {name}
            </li>
          ))}
        </ul>

        <Link to="/pick-target" asChild>
          <Button className="mt-2 w-full hover:cursor-pointer md:h-16 md:text-2xl">
            Começar
          </Button>
        </Link>

        <SocialButtons className="mt-2" />
      </div>
    </div>
  )
}

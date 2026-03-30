import { Button } from "@/components/ui/button"
import { FaDiscord, FaGithub } from "react-icons/fa6"
import { BrowserOpenURL, Quit } from "../../wailsjs/runtime"
import confetti from "canvas-confetti"
import { useEffect } from "react"

export default function FinishedPage() {
  useEffect(() => {
    function randomInRange(min: number, max: number): number {
      return Math.random() * (max - min) + min
    }

    const interval = setInterval(() => {
      confetti({
        angle: randomInRange(80, 100),
        colors: ["#ffbe0b", "#fb5607", "#ff006e", "#8338ec", "#3a86ff"],
        disableForReducedMotion: true,
        drift: randomInRange(-1, 1),
        gravity: randomInRange(0.4, 0.8),
        origin: { x: Math.random(), y: 1.05 },
        particleCount: randomInRange(3, 25),
        scalar: randomInRange(0.4, 1.2),
        startVelocity: randomInRange(50, 100),
        ticks: 3000,
      })
    }, 50)

    const timeout = setTimeout(() => {
      clearInterval(interval)
    }, 3000)

    return () => {
      clearTimeout(timeout)
      clearInterval(interval)
    }
  }, [])

  return (
    <div className="mx-28 flex h-dvh flex-col items-center justify-center select-none">
      <img
        draggable={false}
        className="aspect-square w-2/12 rounded"
        src="cathy.webp"
      />

      <h1 className="mt-5 text-3xl font-bold">Instalação Completa! 🎉</h1>

      <p className="mx-20 mt-2 text-center text-muted-foreground">
        De coração, muito obrigado pela confiança! Foi um esforço de paixão de
        toda a equipe. Caso encontre algum erro ou tenha alguma sugestão, por
        favor, não hesite em nos contar!
        <br />
        No mais, aproveite a viagem. Bom jogo!
      </p>

      <Button
        className="mt-10 h-20 w-full text-2xl hover:cursor-pointer"
        size="lg"
        onClick={Quit}
      >
        Até Lá!
      </Button>
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
            BrowserOpenURL("https://github.com/flyri0/ut-translation-setup-v2")
          }
          className="flex h-auto flex-1 flex-col items-center px-6 py-4 hover:cursor-pointer"
        >
          <div className="flex items-center gap-2 text-lg">
            <FaGithub />
            GitHub
          </div>
          <span className="text-[10px] uppercase opacity-80">Código-fonte</span>
        </Button>
      </div>
    </div>
  )
}

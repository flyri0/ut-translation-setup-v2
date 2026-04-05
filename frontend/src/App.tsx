import { Route, Router, Switch } from "wouter"
import { ThemeProvider } from "./components/theme-provider"
import { TooltipProvider } from "./components/ui/tooltip"
import WelcomePage from "./pages/welcome"
import PickTargetPage from "./pages/pick-target"
import FinishedPage from "./pages/finished"
import InstallPage from "./pages/install"
import { memoryLocation } from "wouter/memory-location"

export default function App() {
  const { hook } = memoryLocation({ path: "/finished" })

  return (
    <ThemeProvider defaultTheme="dark">
      <TooltipProvider>
        <Router hook={hook}>
          <Switch>
            <Route path="/" component={WelcomePage} />
            <Route path="/pick-target" component={PickTargetPage} />
            <Route path="/install" component={InstallPage} />
            <Route path="/finished" component={FinishedPage} />
          </Switch>
        </Router>
      </TooltipProvider>
    </ThemeProvider>
  )
}

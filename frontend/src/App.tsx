import { Route, Router, Switch } from "wouter"
import { ThemeProvider } from "./components/theme-provider"
import { TooltipProvider } from "./components/ui/tooltip"
import WelcomePage from "./pages/welcome"
import PickTargetPage from "./pages/pick-target"

export default function App() {
  return (
    <Router>
      <ThemeProvider defaultTheme="dark">
        <TooltipProvider>
          <Switch>
            <Route path="/" component={WelcomePage} />
            <Route path="/pick-target" component={PickTargetPage} />
          </Switch>
        </TooltipProvider>
      </ThemeProvider>
    </Router>
  )
}

import { Route, Router, Switch } from "wouter"
import { ThemeProvider } from "./components/theme-provider"
import WelcomePage from "./pages/welcome"

export default function App() {
  return (
    <Router>
      <ThemeProvider defaultTheme="dark">
        <Switch>
          <Route path="/" component={WelcomePage} />
        </Switch>
      </ThemeProvider>
    </Router>
  )
}

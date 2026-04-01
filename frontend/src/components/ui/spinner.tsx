import { cn } from "@/lib/utils"
import { FaSpinner } from "react-icons/fa6"

function Spinner({ className, ...props }: React.ComponentProps<"svg">) {
  return (
    <FaSpinner
      role="status"
      aria-label="Loading"
      className={cn("size-4 animate-spin", className)}
      {...props}
    />
  )
}

export { Spinner }

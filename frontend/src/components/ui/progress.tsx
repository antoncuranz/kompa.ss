import * as React from "react"
import { cn } from "@/lib/utils"
import { Primitive } from "@radix-ui/react-primitive";

type PrimitiveDivProps = React.ComponentPropsWithoutRef<typeof Primitive.div>;
interface ProgressProps extends PrimitiveDivProps {
    value?: number | null | undefined;
    secondaryValue?: number | null | undefined;
    max?: number;
    getValueLabel?(value: number, max: number): string;
}

const Progress = ({ className, value, secondaryValue, ...props }: ProgressProps) => (
  <div
    className={cn(
      "relative h-2 w-full overflow-hidden rounded-full bg-secondary",
      className
    )}
    {...props}
  >
    <div
      className="h-full w-full flex-1 transition-all absolute"
      style={{ transform: `translateX(-${100 - (secondaryValue || 0)}%)`, background: "#fdb72e" }}
    />
    <div
      className="h-full w-full flex-1 transition-all absolute"
      style={{ transform: `translateX(-${100 - (value || 0)}%)`, background: "#00175a" }}
    />
  </div>
)

export { Progress }

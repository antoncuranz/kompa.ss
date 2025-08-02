import {FocusEvent, useEffect, useState} from 'react'
import {Input} from "@/components/ui/input.tsx";
import {useToast} from "@/components/ui/use-toast.ts";
import {formatAmount} from "@/components/util.ts";
import {cn} from "@/lib/utils.ts";

interface Props {
  amount: number|null,
  updateAmount: (newAmount: number|null) => void,
  className?: string,
  disabled?: boolean,
  decimals?: number,
  id?: string,
  placeholder?: string,
  warnPredicate?: (value: number|null) => boolean
}

const AmountInput = ({amount, updateAmount, className, disabled = false, decimals = 2, id, placeholder, warnPredicate}: Props) => {
  const [stringAmount, setStringAmount] = useState("")
  const [warning, setWarning] = useState(false)
  const { toast } = useToast();

  useEffect(() => {
    updateAmountInternal()
  }, [amount]);
  function updateAmountInternal() {
    setWarning(amount != null && warnPredicate != null && warnPredicate(amount))
    setStringAmount(formatAmount(amount, decimals))
  }

  function onBlur(event: FocusEvent<HTMLInputElement>) {
    const newAmount = event.target.value

    try {
      const newValue = newAmount ? parseMonetaryValue(newAmount) : null
      updateAmount(newValue)
      setStringAmount(formatAmount(newValue, decimals))
      setWarning(newValue != null && warnPredicate != null && warnPredicate(newValue))
    } catch (e) {
      toast({title: "Unable to parse amount"})
      updateAmountInternal()
    }
  }

  function parseMonetaryValue(valueString: string) {
    const re = new RegExp(String.raw`^(-)?(\d*)(?:[.,](\d{0,${decimals}}))?\d*$`);
    const match = valueString.replace(/\s/g, "").match(re)

    if (!match || match[2].length + (match[3]?.length ?? 0) == 0)
      throw new Error("Unable to parse value string '" + valueString + "'")

    const sign = match[1] ? -1 : 1
    const euros = parseInt(match[2]) || 0
    let cents = parseInt(match[3]) || 0

    if (match[3])
      cents *= Math.pow(10, decimals - match[3].length)

    return sign * (euros * Math.pow(10, decimals!) + cents)
  }

  return (
    <Input id={id} value={stringAmount} placeholder={placeholder} onChange={e => setStringAmount(e.target.value)}
           onBlur={onBlur} className={cn(className, "text-right", warning && "text-yellow-600")} disabled={disabled}/>
  )
}

export default AmountInput
import {Input} from "@/components/ui/input.tsx";
import {useToast} from "@/components/ui/use-toast.ts";
import {Button} from "@/components/ui/button.tsx";
import {Search} from "lucide-react";
import {Location} from "@/types.ts";
import {cn} from "@/lib/utils.ts";

export default function AddressInput({
  address, updateAddress, updateLocation, className
}: {
  address: string,
  updateAddress: (newAddress: string) => void,
  updateLocation: (newAddress: Location|null) => void,
  className?: string,
}) {
  const { toast } = useToast();

  async function searchForAddressUsingMapbox() {
    const mapboxToken = process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN;
    const url = encodeURI("https://api.mapbox.com/search/geocode/v6/forward?q=" + address + "&access_token=" + mapboxToken)

    const response = await fetch(url, {method: "GET"})
    if (response.ok) {
      const json = await response.json()
      const properties = json["features"][0]["properties"]
      updateAddress(properties["full_address"])
      updateLocation({...properties["coordinates"], id: 0})
    } else toast({
      title: "Error finding address. Please try a different format.",
      description: response.statusText
    })
  }

  return (
      <div className={cn("col-span-3 flex gap-2", className)}>
        <Input id="address" value={address}
               onChange={e => updateAddress(e.target.value)}
               className="col-span-3"/>
        <Button type="submit" variant="outline" onClick={() => searchForAddressUsingMapbox()}>
          <Search className="h-3.5 w-3.5"/>
        </Button>
      </div>
  )
}
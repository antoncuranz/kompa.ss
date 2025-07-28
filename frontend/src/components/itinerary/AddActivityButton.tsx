import {Button} from "@/components/ui/button.tsx";

export default function AddActivityButton(){

  return (
    <Button variant="outline" className="m-2 ml-6 mr-6 w-full block border-dashed hover:border-solid" style={{width: "calc(100% - 3rem)"}}>
      Add Activity
    </Button>
  )
}
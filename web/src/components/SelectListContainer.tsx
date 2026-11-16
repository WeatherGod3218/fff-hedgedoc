import { SelectListBlock } from "./SelectListBlock";

export function SelectListContainer() {

    return (
        <div className="mt-4 d-flex flex-column align-items-center w-75 container-sm border mx-0 px-0">
            <SelectListBlock Selected={false} Name="weather" UUID="1234" Icon={true} ></SelectListBlock>
            <SelectListBlock Selected={true} Name="paradox" UUID="1234" Icon={true} ></SelectListBlock>
            <SelectListBlock Selected={true} Name="monkey" UUID="1234" Icon={true} ></SelectListBlock>
            <SelectListBlock Selected={false} Name="kitty" UUID="1234" Icon={true} ></SelectListBlock>
        </div>
    )
}
import { SelectListBlock } from "./SelectListBlock";

export type SelectBlockItem = {
    Name: string;
    UUID: string;
    Type: string;
    Icon: boolean;
    Selected: boolean;
}


type ContainerProps = {
    items: SelectBlockItem[];
    selected: Set<string>;
    onToggle: (uuid: string, type: string) => void;
    invertSelection?: boolean;
};


export function SelectListContainer(props: ContainerProps) {

    return (
        <div className="mt-4 d-flex flex-column align-items-center w-75 container-sm border mx-0 px-0">
            {props.items.map((item) => (
                <SelectListBlock
                    key={item.UUID}
                    Icon={item.Icon}
                    Type={item.Type}
                    Name={item.Name}
                    UUID={item.UUID}
                    Selected={props.invertSelection ? !props.selected.has(item.UUID) : props.selected.has(item.UUID)}
                    onClick={() => props.onToggle(item.UUID, item.Type)}
                />
            ))}
        </div>
    )
}
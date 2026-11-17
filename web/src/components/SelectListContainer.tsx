import { SelectListBlock } from "./SelectListBlock";

type ContainerProps = {
    selected: Set<string>;
    onToggle: (uuid: string) => void;
    invertSelection?: boolean;
};

type User = {
    Name: string;
    UUID: string;
    Icon: boolean
}

export const users: User[] = [
  { Name: "weather", UUID: "1234", Icon: true},
  { Name: "kitty", UUID: "2543", Icon: true},
  { Name: "paradox", UUID: "5623523", Icon: true},
  { Name: "integer", UUID: "12312", Icon: true},
];

export function SelectListContainer(props: ContainerProps) {

    return (
        <div className="mt-4 d-flex flex-column align-items-center w-75 container-sm border mx-0 px-0">
            {users.map((user) => (
                <SelectListBlock
                    key={user.UUID}
                    Icon={user.Icon}
                    Name={user.Name}
                    UUID={user.UUID}
                    Selected={props.invertSelection ? !props.selected.has(user.UUID) : props.selected.has(user.UUID)}
                    onClick={() => props.onToggle(user.UUID)}
                />
            ))}
        </div>
    )
}
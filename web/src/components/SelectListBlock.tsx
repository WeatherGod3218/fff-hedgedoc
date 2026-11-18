import { Check, X } from "lucide-react"


type SelectListBlockProps = {
    Name: string,
    UUID: string,
    Type: string,
    Selected: boolean,
    onClick?: (name: string, uuid: string, selected: boolean) => void,
    Icon?: boolean
}

export function SelectListBlock(props: SelectListBlockProps) {

    return (
        <button type="button" className="btn btn-select-list border-0 border-bottom w-100 py-3 pb-3 d-flex align-items-center gap-2" onClick={() => {
            props.onClick && props.onClick(props.Name, props.UUID, props.Selected)
            }}>
            {props.Icon && (
                <img
                    src={`https://profiles.csh.rit.edu/image/${props.Name}`}
                    alt={`${props.Name}'s profile picture`}
                    width="50"
                    height="50"
                    className="rounded-circle"
                />
            )}
            <h3 className="mb-0">{props.Name}</h3>
            
            <h3 className={`mb-0 ms-auto align-items-center ${props.Selected ? "select-list-approved" : "select-list-denied"}`}>
                {props.Selected ? <>Approved<Check className="mx-1" size={26}/></>  : <>Denied<X className="mx-1" size={26}/></>}
            </h3>
        </button>
    )
}
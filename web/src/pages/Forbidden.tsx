import { Navbar } from "../components/Navbar";
import { useParams } from "react-router";


export function ForbiddenPage() {
    
    const {user} = useParams();

    return (
        <>
            <Navbar/>
            <h1>Sorry {user}, you are not eligable to view this content!</h1>
        </>
    )
}
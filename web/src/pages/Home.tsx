import { Navbar } from "../components/Navbar";

export function HomePage() {
    return (
        <>
            <Navbar/>
            <h1>This is the CSH Home page!</h1>
            <p>The "Notes" page will be where you can view all the notes you have access too. You'll be able to scroll through, and hopefully filter by title</p>
            <p>The "Admin" page will be where you can control access permissions to all of your own personal notes. Specific users and usergroups</p>
            <p>Could I just rewrite the actual source for hedgedoc? Yes. Is this WAY funnier? ABSOLUTELY!</p>
            <h4>GOOD LUCK TO THE FUCKS WHO TRY TO MAINTAIN THIS!</h4>
        </>
    )
}
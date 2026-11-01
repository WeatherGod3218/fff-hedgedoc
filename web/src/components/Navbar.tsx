import { useQuery } from '@tanstack/react-query';

type userInfo = {
    uuid: string,
    username: string,
    full_name: string
}

async function fetchUserInfo() {
  const res = await fetch('/csh/api/v1/user');
  if (!res.ok) throw new Error('Network response error');
  return res.json();
}

export function Navbar() {
    const { data } = useQuery<userInfo>({
        queryKey: ['userInfo'], 
        queryFn: fetchUserInfo,
    });
    
    return (
    <nav className="navbar navbar-expand-lg bg-primary sticky-top shadow-sm" data-bs-theme="dark">
        <div className="container-lg">
        <a href="/csh" className="navbar-brand d-flex align-items-center py-0 fs-4 text-light">
            <img className="object-fit-contain ms-2" height="32" src="https://assets.csh.rit.edu/pubsite/csh_logo_square.svg"/>
        </a>
        <a href='/' className="navbar-brand d-flex align-items-center py-0 fs-4 text-light">
            <img className="object-fit-contain" height="32" src="/banner/banner_h_wb.svg"/>
        </a>

        <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbar-lg-1" aria-controls="navbar-lg-1" aria-expanded="false" aria-label="Toggle navigation">
            <span className="navbar-toggler-icon"></span>
        </button>
        
        <div className="collapse navbar-collapse text-light " id="navbar-lg-1">
            <ul className="navbar-nav align-items-center">
            <li className="nav-item">
                <a href="/csh/notes" className="nav-link pe-2 py-0 link-light">Notes</a>
            </li>
            <li className="nav-item">
                <a href="/csh/admin" className="nav-link px-2 py-0 link-light">Admin</a>
            </li>
            </ul>
            <div className="bd-theme dropdown ms-auto">
                <button className="btn dropdown-toggle text-body" data-bs-toggle="dropdown" role="button" aria-haspopup="true" aria-expanded="false">
                    <i className="theme-icon-active bi bi-circle-half"></i>
                    <span className="bd-theme-text"></span>
                </button>
                <ul className="dropdown-menu dropdown-menu-end" aria-labelledby="bd-theme">
                    <li><button className="dropdown-item" aria-pressed="false" data-bs-theme-value="light">
                        <i className="bi bi-sun-fill pe-2"></i>
                        Light
                        </button>
                    </li>
                    <li>
                        <button className="dropdown-item" aria-pressed="false" data-bs-theme-value="dark">
                        <i className="bi bi-moon-fill pe-2"></i>
                        Dark
                        </button>
                    </li>
                    <li>
                        <button className="dropdown-item" aria-pressed="false" data-bs-theme-value="auto">
                        <i className="bi bi-circle-half pe-2"></i>
                        Auto
                        </button>
                    </li>
                </ul>
            </div>
            <div className="dropdown text-end">
            <div className="">
                <a
                    href="#"
                    className="d-flex align-items-center link-body-emphasis text-decoration-none dropdown-toggle"
                    data-bs-toggle="dropdown"
                    aria-expanded="false"
                >
                <p className="px-2 mb-0">{data ? data.full_name : "Loading"}</p>
                <img
                    src={data ? "https://profiles.csh.rit.edu/image/" + data.username : "https://assets.csh.rit.edu/pubsite/csh_logo_square.svg"}
                    alt="profile picture"
                    width="32"
                    height="32"
                    className="rounded-circle"
                    />
                </a>
                <ul className="dropdown-menu dropdown-menu-end text-small">
                    <li><a className="dropdown-item" href="/new">New Note</a></li>
                    <li><a className="dropdown-item" href={data ? "http://profiles.csh.rit.edu/user/" + data.username : "http://profiles.csh.rit.edu"}>Profile</a></li>
                    <li><hr className="dropdown-divider" /></li>
                    <li><a className="dropdown-item" href="/auth/logout">Sign out</a></li>
                </ul>
                </div>
    
            </div>
        </div>
        </div>
    </nav>
  )
}
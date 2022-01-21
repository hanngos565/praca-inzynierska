import 'bulma/css/bulma.min.css';
import '@fortawesome/fontawesome-free/css/all.min.css';

const burgerIconIsActive = () => {
    const navbarMenu = document.querySelector("#navbar-menu")
    navbarMenu.classList.toggle('is-active')
}

const Header = () => {
    return (
        <nav className="navbar is-transparent has-shadow">
            <div className="navbar-brand">
                <div className="navbar-burger" id="navbar-burger" onClick={burgerIconIsActive}>
                    <span></span>
                    <span></span>
                    <span></span>
                </div>
            </div>

            <div id="navbar-menu" className="navbar-menu">
                <div className="navbar-start">
                    <a className="navbar-item" href={"/"}>
                        Home
                    </a>
                    <a className="navbar-item " href={"/#configuration"}>
                        Configuration
                    </a>
                    <a className="navbar-item " href={"/#simulation"}>
                        Simulation
                    </a>
                    <div className="navbar-item has-dropdown is-hoverable">
                        <a className="navbar-link" href={"/#results"}>
                            Results
                        </a>
                        <div className="navbar-dropdown is-boxed">
                            <a className="navbar-item" href={"/simulation-results/demo"}>
                                Demo
                            </a>
                        </div>
                    </div>
                </div>
                <div className="navbar-end">
                    <div className="navbar-item">
                        <div className="field is-grouped">
                            <p className="control">
                                <a className="button" target="_blank"
                                   href="https://github.com/hanngos">
                                <span className="icon">
                                    <i className="fab fa-github"/>
                                </span>
                                    <span>
                                    GitHub
                                </span>
                                </a>
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </nav>)
}

export {Header}
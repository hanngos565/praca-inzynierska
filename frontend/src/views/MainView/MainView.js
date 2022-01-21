import React from "react";
import {Header} from "../../components/Header/Header";
import {Footer} from "../../components/Footer/Footer";
import {Configuration} from "./ConfigurationPage/Configuration";
import {Simulation} from "./SimulationPage/Simulation";
import {Results} from "./ResultsPage/Results";

const MainView = () => {
    return (
        <section className="hero is-fullheight">
            <div className="hero-head">
                <Header/>
            </div>

            <div className="hero-body">
                <div className="container">
                    <Configuration/>
                    <Simulation/>
                    <Results/>
                </div>
            </div>

            <div className="hero-foot">
                <Footer/>
            </div>
        </section>
    )
}
export {MainView}
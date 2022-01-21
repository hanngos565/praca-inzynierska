import React from "react";
import {getResults} from "../../API";
import {algorithmID} from "../utils";
import "bulma-extensions/bulma-accordion/dist/css/bulma-accordion.min.css"
import {Header} from "../../components/Header/Header";
import {Footer} from "../../components/Footer/Footer";
import {DemoParameters, SimulationResults} from "./components/DemoResults";

export const dict = { "in-progress" : "is-link" ,
    "finished" : "is-primary" ,
    "error" : "is-danger" ,
};

const ResultsView = () => {
    const [results, setResults] = React.useState()

    const handleClick = () => {
        let acc = document.getElementsByClassName("accordion");
        for (let i = 0; i < acc.length; i++) {
            acc[i].addEventListener("click", function () {
                if (this.classList.contains('is-active'))
                    this.classList.remove('is-active')
                else this.classList.add('is-active')
            });
        }
    }
    setInterval(() => getResults("demo", algorithmID), 1000);

    React.useEffect(() => {
        const fetchData = async () => {
            let r = await getResults("demo", algorithmID);
            r = r?.sort((ts1, ts2) => new Date(ts2.timeStamp) - new Date(ts1.timeStamp));
            setResults(r);
        }
        fetchData()
            .catch(console.error);
    }, [])

    return (
        <section className="hero is-fullheight">
            <div className="hero-head">
                <Header/>
            </div>
            <div className="hero-body">
                <div className="container">
                    <section className="section notification">
                        <h1 className="title">Demo Results</h1>
                        <section className="accordions">
                            {results &&
                                results.map(r => {
                                    return (
                                        <div key={r.timeStamp} className={"accordion " + dict[r.status]}>
                                            <div className="accordion-header toggle" onClick={handleClick}>
                                                <p>{new Date(r.timeStamp).toLocaleDateString("en-US", {
                                                    year: "2-digit",
                                                    month: "2-digit",
                                                    day: "2-digit",
                                                    hour12: false,
                                                    hour: "2-digit",
                                                    minute: "2-digit",
                                                    second: "2-digit"
                                                })}</p>
                                                <button className="toggle" aria-label="toggle"/>
                                            </div>
                                            <div className="accordion-body">
                                                <div className="container accordion-content">
                                                    <div className="tile is-ancestor">
                                                        <DemoParameters r={r}/>
                                                        {r.result && <SimulationResults r={r}/>}
                                                    </div>
                                                </div>
                                            </div>
                                        </div>)
                                })}
                        </section>
                    </section>
                </div>
            </div>
            <div className="hero-foot">
                <Footer/>
            </div>
        </section>
    )
}
export {ResultsView}
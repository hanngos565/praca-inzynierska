import React from "react";

const Results = () => {
    return (
        <section id={"results"} className="section notification is-info">
            <h1 className="title">Results</h1>
            <div className="columns">
                <div className="column">
                    <a href={"/simulation-results/demo"}>
                    <button id={"demo"} className="button is-fullwidth is-medium">
                        <p>Get<strong> Demo </strong>Results</p>
                    </button>
                    </a>
                </div>
                <div className="column">
                    <button id={"..."} className="button is-fullwidth is-medium" disabled={true}>
                        <p>Get<strong> ... </strong>Results</p>
                    </button>
                </div>
            </div>
        </section>
    )
}

export {Results}
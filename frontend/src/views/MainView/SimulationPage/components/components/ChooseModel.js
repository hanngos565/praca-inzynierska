import React from "react";
import {getModels} from "../../../../../API";
import {algorithmID} from "../../../../utils";

const ChooseModel = () => {
    const [models, setModels] = React.useState();

    React.useEffect(() => {
        const fetchData = async () => {
            const m = await getModels(algorithmID);
            setModels(m);
        }
        fetchData()
            .catch(console.error);
    }, [])

    return (
        <div>
            <p className="subtitle">Choose Model</p>
            <div className="container has-text-centered">
                <div className="select is-link ">
                    <select id="choose-model">
                        {models &&
                            models.map(m => {
                                return <option key={m}>{m}</option>;
                            })}
                    </select>
                </div>
            </div>
        </div>
    )
}

export {ChooseModel}
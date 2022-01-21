import {dict} from "../ResultsView";
import React from "react";

const DemoParameters = (props) => {
    return (
        <div className="tile is-parent">
            <div className="tile is-child box">
                <p className="subtitle">Simulation Parameters</p>
                <div className="container has-text-centered">
                    <table className="table is-bordered table is-striped">
                        <thead>
                        <tr>
                            <td>model</td>
                            <td>{props.r.model}</td>
                        </tr>
                        <tr>
                            <td>algorithm</td>
                            <td>{props.r.algorithm}</td>
                        </tr>
                        <tr>
                            <td>image</td>
                            <td><img id={props.r.timeStamp} src={props.r.image} alt={props.r.timeStamp}/></td>
                        </tr>
                        </thead>
                    </table>
                </div>
            </div>
        </div>
    )
}

const SimulationResults = (props) => {
    const canvasWidth = window.innerWidth * .5;
    const canvasHeight = window.innerHeight * .5;

    const drawImageActualSize = (canvas, ctx, image) => {
        canvas.width = image.naturalWidth;
        canvas.height = image.naturalHeight;
        ctx.drawImage(image, 0, 0);
    }

    function draw(ctx, result) {
        for (let i = 0; i < result.names.length; i++) {
            let bbox = result.bbox[i];
            ctx.strokeStyle = 'red';
            ctx.lineWidth = 2;
            ctx.strokeRect(bbox[1], bbox[0], bbox[3] - bbox[1], bbox[2] - bbox[0]);
            ctx.fillText(result.names[i] + " " + result.scores[i].toString().slice(0, 6), bbox[1], bbox[0])
            ctx.restore();
        }
    }

    React.useEffect(() => {
            const fetchData = async () => {
                const canvas = document.getElementById("canvas" + props.r.timeStamp)
                const ctx = canvas.getContext('2d');
                const image = new Image(document.getElementById(props.r.timeStamp).width, document.getElementById(props.r.timeStamp).height)
                image.src = props.r.image
                drawImageActualSize(canvas, ctx, image);
                if (props.r.result !== '' && props.r.result !== `"error"`) {
                    draw(ctx, JSON.parse(props.r.result))
                }
            }
            fetchData()
                .catch(console.error);
    }, []);

    const downloadImage = (e) => {
        const type = "png"
        let canvas = document.getElementById("canvas" + e.target.id);

        let anchor = document.createElement("a");
        anchor.download = "download." + type;
        anchor.href = canvas.toDataURL("image/" + type);

        anchor.click();
        anchor.remove();
    }

    return(
        <div className="tile is-parent">
            <div className="tile is-child box">
                <p className="subtitle">Simulation Results</p>
                <div className="container has-text-centered">
                    <canvas
                        className="App-canvas"
                        id={"canvas" + props.r.timeStamp}
                        width={canvasWidth}
                        height={canvasHeight}
                    />
                    <button id={props.r.timeStamp} className={"button "+ dict[props.r.status]} onClick={downloadImage}>Download</button>
                </div>
            </div>
        </div>
    )
}

export {DemoParameters, SimulationResults}
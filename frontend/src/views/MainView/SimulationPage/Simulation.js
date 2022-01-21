import React from "react";
import {DemoSection} from "./components/DemoSection";
import {ImagesModal} from "./components/ImagesModal";

const Simulation = () => {
    const tabs = document.querySelectorAll('.tabs li')
    const tabContentBoxes = document.querySelectorAll('#tab-content > div')
    tabs.forEach((tab) => {
        tab.addEventListener('click', () => {
            tabs.forEach(item => item.classList.remove('is-active'))
            tab.classList.add('is-active')

            const target = tab.dataset.target
            tabContentBoxes.forEach(box => {
                if (box.getAttribute('id') === target) {
                    box.classList.remove('is-hidden')
                } else {
                    box.classList.add('is-hidden')
                }
            })
        })
    })
    return (
        <section id={"simulation"} className="section notification is-link is-light">
            <h1 className="title">Simulation</h1>
            <div className="tabs is-medium is-toggle is-boxed">
                <ul>
                    <li className="is-active" data-target={"DEMO"}>
                        <a>
                            Demo
                        </a>
                    </li>
                    <li className="is-disabled" data-target={"..."}>
                        <a>
                            ...
                        </a>
                    </li>
                </ul>
            </div>
            <div id={"tab-content"}>
                <div id={"DEMO"} className="notification is-link">
                    <DemoSection/>
                </div>
                <div id={"..."} className="notification is-link is-hidden is-disabled">
                    <p>...</p>
                </div>
            </div>
            <ImagesModal/>
        </section>
    )
}

export {Simulation}
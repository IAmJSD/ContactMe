import React from "react"
import ReactDOM from "react-dom"

let config

class FormSelector extends React.Component {
    constructor(props) {
        super(props)
    }

    render() {
        return <div>
            <h5 className="title is-5">What is the form you wish to select?</h5>
            <div className="select">
                <select onInput={e => this.props.formSelected(e.target.value)}>
                    <option value="" selected disabled hidden>Form selection required</option>
                    {
                        Object.keys(config.forms).map(x => <option>{x}</option>)
                    }
                </select>
            </div>
            <hr />
        </div>
    }
}

class FormOption extends React.Component {
    render() {
        return <div>
            <h5 className="title is-5">{this.props.name} {this.props.value.required ? "*" : ""}</h5>
            <hr />
        </div>
    }
}

class Form extends React.Component {
    render() {
        return <div>
            {this.props.form.warning ? <div className="notification is-warning">{this.props.form.warning}</div> : null}
            {Object.keys(this.props.form.children).map(x => <FormOption onChange={result => this.props.properties[x] = result} name={x} value={this.props.form.children[x]} />)}
        </div>
    }
}

class App extends React.Component {
    constructor(props) {
        super(props)
        this.state = {properties: {}}
    }

    render() {
        return <div className="container">
            <div style={{textAlign: "center", paddingTop: 20}}>
                <h1 className="title is-1">{ config.page_description.title }</h1>
                <h5 className="title is-5">{ config.page_description.description }</h5>
            </div>
            <hr />
            <FormSelector formSelected={selectedForm => this.setState({selectedForm, form: config.forms[selectedForm], properties: {}})} />
            {
                this.state.form ? this.state.form.type === "error" ? <div className="notification is-warning">{this.state.form.message}</div> : <Form
                    properties={this.state.properties}
                    name={this.state.selectedForm}
                    form={this.state.form}
                /> : null
            }
        </div>
    }
}

fetch("/_forms").then(async res => {
    config = await res.json()
    ReactDOM.render(<App />, document.getElementById("app"))
})

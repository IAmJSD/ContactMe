import React from "react"
import ReactDOM from "react-dom"

// Handles the config.
let config

// Contains all of the values.
const values = []

// Allow the user to select a form.
class FormSelector extends React.Component {
    constructor(props) {
        super(props)
        this.state = {initialOption: <option value="">Form selection required</option>}
    }

    formSelected(e) {
        this.setState({initialOption: null})
        for (const form of config.forms) {
            if (form.name === e.target.value) this.props.formSelected(form)
        }
    }

    render() {
        return <div>
            <h5 className="title is-5">What is the form you wish to select?</h5>
            <div className="select">
                <select onInput={e => this.formSelected(e)}>
                    {this.state.initialOption}
                    {config.forms.map(x => <option key={x.name}>{x.name}</option>)}
                </select>
            </div>
            <hr />
        </div>
    }
}

// Handles a e-mail input.
class EmailInput extends React.Component {
    render() {
        return <div className="field-body">
            <div className="field">
                <p className="control">
                    <input className="input" type="email" onChange={e => this.props.onChange(e.target.value === "" ? undefined : e.target.value)} />
                </p>
            </div>
        </div>
    }
}

// Handles a small textbox input.
class SmallTextboxInput extends React.Component {
    render() {
        return <div className="field-body">
            <div className="field">
                <p className="control">
                    <input className="input" type="text" onChange={e => this.props.onChange(e.target.value === "" ? undefined : e.target.value)} />
                </p>
            </div>
        </div>
    }
}

// Handles a large textbox input.
class LargeTextboxInput extends React.Component {
    render() {
        return <div className="field-body">
            <div className="field">
                <textarea className="textarea" onChange={e => this.props.onChange(e.target.value === "" ? undefined : e.target.value)} rows="10" />
            </div>
        </div>
    }
}

// Handles the integer input.
class IntInput extends React.Component {
    render() {
        return <div className="field-body">
            <div className="field">
                <p className="control">
                    <input className="input" type="number" onChange={e => this.props.onChange(e.target.value === "" ? undefined : Number(e.target.value) || undefined)} />
                </p>
            </div>
        </div>
    }
}

// Handles the form option.
class FormOption extends React.Component {
    constructor(props) {
        super(props)
        this.option = this.props.value.type === "email" ? <EmailInput onChange={result => this.props.onChange(result)} config={this.props.value} /> :
            this.props.value.type === "int" ? <IntInput onChange={result => this.props.onChange(result)} config={this.props.value} /> :
            this.props.value.type === "small_textbox" ? <SmallTextboxInput onChange={result => this.props.onChange(result)} config={this.props.value} /> :
            this.props.value.type === "large_textbox" ? <LargeTextboxInput onChange={result => this.props.onChange(result)} config={this.props.value} /> :
            null
    }

    render() {
        return <div>
            <h5 className="title is-5">{this.props.name} <span style={{color: "red"}}>{this.props.value.required ? "*" : ""}</span></h5>
            {this.option}
            <hr />
        </div>
    }
}

// Handles the base for the form.
class Form extends React.Component {
    render() {
        return <div>
            {this.props.form.warning ? <div className="notification is-warning">{this.props.form.warning}</div> : null}
            {Object.keys(this.props.form.children).map(x => <FormOption key={x} onChange={result => this.props.properties[x] = result} name={x} value={this.props.form.children[x]} />)}
        </div>
    }
}

// Shows the submit button.
class SubmitButton extends React.Component {
    render() {
        return <div className="field">
            <p className="control">
                <button className="button is-link" onClick={() => this.props.onSubmit()}>
                    Submit
                </button>
            </p>
        </div>
    }
}

// Handles the main app.
class App extends React.Component {
    constructor(props) {
        super(props)
        this.state = {properties: {}}
    }

    formSelected(form) {
        this.setState({selectedForm: form.name, form, properties: {}})
        values.length = 0
        for (const key in form.children || {}) values.push(key)
    }

    submitButton() {
        // Create the object.
        const obj = {__formName: this.state.selectedForm}
        for (const field of values) {
            const properties = this.state.form.children[field]
            const value = this.state.properties[field]
            if (properties.required && !value) {
                alert("Required fields are missing.")
                return
            }
            if (value) obj[field] = String(value)
        }

        // POST the object.
        this.postObject(obj)
    }

    postObject(obj) {
        fetch("/", {method: "POST", body: JSON.stringify(obj)}).catch(err => alert(err)).then(res => {
            if (!res.ok) {
                alert("Form submission failed.")
                return
            }
            alert("Form submitted.")
            window.location.replace(config.redirect)
        })
    }

    render() {
        return <div className="container">
            <div style={{textAlign: "center", paddingTop: 20}}>
                <h1 className="title is-1">{ config.page_description.title }</h1>
                <p>{ config.page_description.description }</p>
            </div>
            <hr />
            <FormSelector formSelected={form => this.formSelected(form)} />
            {
                this.state.form ? this.state.form.type === "error" ? <div className="notification is-warning">{this.state.form.message}</div> : <Form
                    properties={this.state.properties}
                    name={this.state.selectedForm}
                    form={this.state.form}
                /> : null
            }
            {(this.state.form || {type: "error"}).type !== "error" ? <SubmitButton onSubmit={() => this.submitButton()} /> : null}
        </div>
    }
}

// Loads the page.
fetch("/_forms").then(async res => {
    config = await res.json()
    ReactDOM.render(<App />, document.getElementById("app"))
})

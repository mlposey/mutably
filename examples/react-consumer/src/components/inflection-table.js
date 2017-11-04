import React, { Component } from 'react';
import '../styles/inflection-table.css';

/** Renders an inflection/conjugation table for a verb */
class InflectionTable extends Component {
    render() {
        return (
            <table><tbody>
                <tr>
                    <th colSpan="3"
                        style={{textAlign: 'center'}}>
                        {this.props.inf.Infinitive}
                    </th>
                </tr>
                <tr>
                    <th />
                    <th>Present Tense</th>
                    <th>Past Tense</th>
                </tr>
                <tr>
                    <th>1st person singular</th>
                    <th>{this.props.inf.Present.First.join(', ')}</th>
                    <th>{this.props.inf.Past.First.join(', ')}</th>
                </tr>
                <tr>
                    <th>2nd person singular</th>
                    <th>{this.props.inf.Present.Second.join(', ')}</th>
                    <th>{this.props.inf.Past.Second.join(', ')}</th>
                </tr>
                <tr>
                    <th>3rd person singular</th>
                    <th>{this.props.inf.Present.Third.join(', ')}</th>
                    <th>{this.props.inf.Past.Third.join(', ')}</th>
                </tr>
                <tr>
                    <th>Plural</th>
                    <th>{this.props.inf.Present.Plural.join(', ')}</th>
                    <th>{this.props.inf.Past.Plural.join(', ')}</th>
                </tr>
            </tbody></table>
        );
    }
}

export default InflectionTable;
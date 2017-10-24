module Msg exposing (..)

-- Core packages

import Http


-- Local packages

import Mode.Scoreboard exposing (Scoreboard)
import Mode.Message exposing (Message)
import Mode.SubmitForm as SubmitForm exposing (SubmitResponse)


type ViewMode
    = ScoreboardView
    | SubmitForm
    | MessagesView


type Msg
    = SwitchMode ViewMode
    | ScoreboardRetrieved (Result Http.Error Scoreboard)
    | FlagSubmitted (Result Http.Error SubmitResponse)
    | MessagesRetrieved (Result Http.Error (List Message))
    | GotInput SubmitForm.Input

module Main exposing (..)

-- Core imports

import Html exposing (..)
import Html.Attributes exposing (..)


-- Package imports

import Http


-- Local imports

import Action exposing (Action(..), Team, Scoreboard, fetchScoreboard, viewScoreboard)


-- MAIN


main : Platform.Program Basics.Never Model Msg
main =
    Html.program
        { init = init
        , update = update
        , view = view
        , subscriptions = subscriptions
        }



-- MODEL


type alias Model =
    { teams : List Team
    , inProgress : Action
    }


init : ( Model, Cmd Msg )
init =
    ( Model [] Action.ViewScoreboard, Http.send GotScoreboard fetchScoreboard )



-- UPDATE


type Msg
    = GotScoreboard (Result Http.Error Scoreboard)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GotScoreboard (Ok scoreboard) ->
            ( Model scoreboard.teams Action.Nothing, Cmd.none )

        GotScoreboard (Err _) ->
            -- ( model, Http.send GotScoreboard fetchScoreboard )
            ( model, Cmd.none )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- VIEW


view : Model -> Html Msg
view model =
    div []
        [ viewNav model
        , viewHeading model
        , viewScoreboard { teams = model.teams }
        ]


viewHeading : Model -> Html Msg
viewHeading model =
    div [ id "heading" ]
        [ h1 [] [ text "CTF Scoreboard" ]
        ]


viewNav : Model -> Html Msg
viewNav model =
    div [ id "navigation" ]
        [ a [ href "/login" ] [ text "Admin login" ]
        , span [] [ text " | " ]
        , a [ href "/submit" ] [ text "Submit a flag" ]
        , span [] [ text " | " ]
        , a [ href "/message" ] [ text "Messages" ]
        ]

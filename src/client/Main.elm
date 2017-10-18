module Main exposing (main)

-- Standard imports

import Html exposing (..)
import Html.Attributes exposing (..)


-- Local imports

import Mode.Scoreboard as Scoreboard exposing (Scoreboard(..))


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


type alias Date =
    String


type ViewMode
    = ScoreboardView Scoreboard
    | FormView
    | MessagesView


type alias Model =
    { mode : ViewMode
    }


init : ( Model, Cmd Msg )
init =
    ( Model (ScoreboardView testScoreboard), Cmd.none )


testScoreboard : Scoreboard
testScoreboard =
    Scoreboard
        [ { rank = 1, name = "Team one", score = 150, lastSubmission = "9:30" }
        , { rank = 2, name = "H4xx0R", score = 132, lastSubmission = "9:25" }
        , { rank = 3, name = "31337", score = 130, lastSubmission = "9:46" }
        , { rank = 4, name = "b4d455", score = 80, lastSubmission = "9:15" }
        , { rank = 5, name = "CTF TO", score = 68, lastSubmission = "9:08" }
        , { rank = 6, name = "T.", score = 35, lastSubmission = "8:12" }
        , { rank = 7, name = "DC416", score = 10, lastSubmission = "7:50" }
        ]



-- UPDATE


type Msg
    = SwitchMode ViewMode


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        SwitchMode mode ->
            ( { model | mode = mode }, Cmd.none )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- VIEW


view : Model -> Html Msg
view model =
    div []
        [ ul [ class "horizontal buttonList" ]
            [ li [] [ a [ class "btn btnPrimary waves-effect waves-light" ] [ text "Submit" ] ]
            , li [] [ a [ class "btn btnSecondary waves-effect waves-light" ] [ text "Messages" ] ]
            , li [ style [ ( "float", "right" ) ] ] [ a [ class "btn btnSecondary waves-effect waves-light" ] [ text "Admin" ] ]
            ]
        , viewMode model.mode
        ]


viewMode : ViewMode -> Html Msg
viewMode mode =
    div [] <|
        case mode of
            ScoreboardView scoreboard ->
                [ Scoreboard.view scoreboard ]

            _ ->
                []
